package handlers

import (
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LanguageHandler struct {
	Store *stores.LanguageStore
}

func (app LanguageHandler) GetLanguages(w http.ResponseWriter, r *http.Request) {
	languages := app.Store.GetLanguages()
	json, err := json.Marshal(languages)

	if err != nil {
		log.Error(err)
	}

	fmt.Fprintf(w, string(json))
}
