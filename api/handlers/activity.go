package handlers

import (
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ActivityHandler struct {
	Store *stores.ActivityStore
}

func (app ActivityHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	activities := app.Store.GetActivities()
	json, err := json.Marshal(activities)

	if err != nil {
		log.Error(err)
	}

	fmt.Fprintf(w, string(json))
}
