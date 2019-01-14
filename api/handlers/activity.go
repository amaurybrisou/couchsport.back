package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/stores"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type activityHandler struct {
	Stores *stores.StoreFactory
}

//All returns all the activities in DB
func (app activityHandler) All(w http.ResponseWriter, r *http.Request) {
	activities, err := app.Stores.ActivityStore().All()
	if err != nil {
		log.Error(err)
		http.Error(w, "error retrieving activities", http.StatusInternalServerError)
	}

	json, err := json.Marshal(activities)
	if err != nil {
		log.Error(err)
		http.Error(w, "error retrieving activities", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, string(json))
}
