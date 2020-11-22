package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/goland-amaurybrisou/couchsport/config"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const prefix = "/api"

type Instance struct {
	Db         *gorm.DB
	C          *config.Config
	HTTPServer *http.Server
	router     *http.ServeMux
}

var s *Instance

//NewInstance creates the app server object
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
	go func() {
		log.Printf("Listenning on  %s:%s", s.C.Listen, strconv.Itoa(s.C.Port))
		if err := s.HTTPServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

//Shutdown http server
func (s *Instance) Shutdown() {
	if s.Db != nil {
		s.Db.Close()
	}

	if s.HTTPServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := s.HTTPServer.Shutdown(ctx)
		if err != nil {
			cancel()
			log.Panic(err)
		} else {
			cancel()
			s.HTTPServer = nil
		}
	}
}

// RegisterHandler registers api handlers
func (s *Instance) RegisterHandler(path string, handler http.HandlerFunc) {
	log.Infof("registering handler at path in %s environment %s, cors is enabled in dev", prefix+path, s.C.Env)
	if s.C.Env == "dev" {
		handler = s.enableCors(handler)
	}
	s.router.Handle(prefix+path, handler)
}

//ServerPublic ...
func (s *Instance) ServePublic(path string) {
	log.Infof("serving files at %s", http.Dir(path))
	s.router.Handle("/static/", http.StripPrefix(`/static/`, http.FileServer(http.Dir(path+"/static"))))
	s.router.Handle("/lib/", http.StripPrefix(`/lib/`, http.FileServer(http.Dir(path+"/lib"))))
	s.router.Handle("/uploads/", http.StripPrefix(`/uploads/`, http.FileServer(http.Dir(path+"/uploads"))))
	s.router.Handle("/fonts/", http.StripPrefix(`/fonts/`, http.FileServer(http.Dir(path+"/fonts"))))
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path+"/index.html")
	})
}

func (s *Instance) enableCors(pass http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:"+strconv.Itoa(s.C.Port+1))
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+strconv.Itoa(s.C.Port+1))
		w.Header().Set("Access-Control-Request-Headers", "X-Requested-With")
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		pass(w, r)
	}
}
