package server

import (
	"context"
	"couchsport/config"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Instance struct {
	Db         *gorm.DB
	C          *config.Config
	HTTPServer *http.Server
	router     *http.ServeMux
}

var s *Instance

func NewInstance(c *config.Config) *Instance {
	if s != nil {
		return s
	}
	r := http.NewServeMux()

	db := mustOpenDb(c)
	// Setup(s.C, s.Db)

	s = &Instance{
		C:          c,
		router:     r,
		HTTPServer: &http.Server{Addr: c.Listen + ":" + strconv.Itoa(c.Port), Handler: r},
		Db:         db,
		// just in case you need some setup here
	}

	return s
}

//Start the current Instance
func (s *Instance) Start() {

	defer s.Db.Close()

	signalChan := make(chan os.Signal, 1)
	signalDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		log.Printf("Listenning on  %s:%s", s.C.Listen, strconv.Itoa(s.C.Port))
		if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	go func() {
		<-signalChan
		log.Info("received os.Interrupt signal, stopping services")
		if err := s.HTTPServer.Shutdown(nil); err != nil {
			log.Panic(err)
		}
		close(signalDone)
	}()

	<-signalDone

}
func (s *Instance) Shutdown() {
	if s.HTTPServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := s.HTTPServer.Shutdown(ctx)
		if err != nil {
			cancel()
			log.Panic(err)
		} else {
			s.HTTPServer = nil
		}
	}
}

const prefix = "/api"

func (s *Instance) RegisterHandler(path string, handler http.HandlerFunc) {
	log.Info("registering handler at path ", prefix+path)
	s.router.Handle(prefix+path, handler)
}

func (s *Instance) ServePublic(path string) {
	log.Infof("serving files at %s", http.Dir(path))
	s.router.Handle(`/static/`, http.FileServer(http.Dir(path)))
	s.router.Handle(`/uploads/`, http.FileServer(http.Dir(path)))
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path+"index.html")
	})
}
