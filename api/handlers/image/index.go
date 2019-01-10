package image

import (
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ImageHandler struct {
	Store stores.ImageStore
}

func (app ImageHandler) SoftDeleteHandler(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	result, err := app.Store.Delete(userID, r.Body)
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
