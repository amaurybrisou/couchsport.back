package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/stores"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type languageHandler struct {
	Store *stores.StoreFactory
}

//All returns all languages
func (app languageHandler) All(w http.ResponseWriter, r *http.Request) {
	languages, err := app.Store.LanguageStore().All()
	if err != nil {
		log.Error(err)
		http.Error(w, http.ErrNotSupported.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(languages)

	if err != nil {
		log.Error(err)
		http.Error(w, http.ErrNotSupported.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(json))
}
