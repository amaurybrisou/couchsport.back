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

// func (app ProfileHandler) UploadAvatarHandler(userID uint, w http.ResponseWriter, r *http.Request) {
// 	r.Close = true

// 	if r.Body != nil {
// 		defer r.Body.Close()
// 	}

// 	imagePath, err := app.FileStore.FileUpload(r, userID)
// 	//here we call the function we made to get the image and save it
// 	if err != nil {
// 		http.Error(w, fmt.Errorf("could not upload image %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, fmt.Errorf("could not parse form %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	user, err := app.UserStore.GetUserByID(userID)
// 	if err != nil {
// 		http.Error(w, fmt.Errorf("could not find user %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	user.Profile.Avatar = imagePath

// 	if errs := app.Store.Db.Save(&user.Profile).GetErrors(); len(errs) > 0 {
// 		for _, err := range errs {
// 			log.Error(err)
// 		}
// 		http.Error(w, fmt.Errorf("error updating user %s", errs[0]).Error(), http.StatusBadRequest)
// 	}

// 	json, err := json.Marshal(user.Profile)

// 	if err != nil {
// 		http.Error(w, fmt.Errorf("could not encode output %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Fprintf(w, string(json))

// }
