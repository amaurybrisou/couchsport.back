package page

import (
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PageHandler struct {
	Store *stores.PageStore
}

func (app PageHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	pages := app.Store.GetPages()
	json, err := json.Marshal(pages)

	if err != nil {
		log.Error(err)
	}

	fmt.Fprintf(w, string(json))

}

func (app PageHandler) CreateHandler(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	pageObj, err := app.Store.CreateOrUpdate(userID, r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("could not create page %s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(pageObj)

	if err != nil {
		log.Error(err)
	}

	fmt.Fprint(w, string(json))

}

func (app PageHandler) DeleteHandler(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	result, err := app.Store.Delete(r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("could not extract body %s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(struct{ Result bool }{Result: result})

	if err != nil {
		log.Error(err)
	}

	fmt.Fprint(w, string(json))
}

func (app PageHandler) PublishHandler(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	result, err := app.Store.Publish(userID, r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("could not extract body %s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(struct{ Result bool }{Result: result})

	if err != nil {
		log.Error(err)
	}

	fmt.Fprint(w, string(json))
}
