package profile

import (
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	// log "github.com/sirupsen/logrus"
	"net/http"
)

type ProfileHandler struct {
	UserStore stores.UserStore
	Store     stores.ProfileStore
	// FileStore *stores.FileStore
}

func (app ProfileHandler) GetProfileHandler(UserId uint, w http.ResponseWriter, r *http.Request) {
	profile, err := app.Store.GetProfileByOwnerID(UserId)

	if err != nil {
		http.Error(w, fmt.Errorf("could not find User %s", err).Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(profile)

	if err != nil {
		http.Error(w, fmt.Errorf("could encode output %s", err).Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(json))

}

func (app ProfileHandler) UpdateProfile(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	profile, err := app.Store.Update(userID, r.Body)
	if err != nil {
		http.Error(w, fmt.Errorf("could not update profile").Error(), http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(profile)

	if err != nil {
		http.Error(w, fmt.Errorf("could encode output %s", err).Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(json))

}
