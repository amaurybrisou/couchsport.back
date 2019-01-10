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

// func (app ImageHandler) UploadHandler(UserID uint, w http.ResponseWriter, r *http.Request) {

// 	imagePath, err := app.FileStore.FileUpload(r, UserID)
// 	//here we call the function we made to get the image and save it
// 	if err != nil {
// 		http.Error(w, fmt.Errorf("could not upload image %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, fmt.Errorf("could not parse form %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	pageID := r.PostFormValue("PageID")

// 	id, err := strconv.ParseUint(pageID, 10, 32)
// 	if err != nil {
// 		http.Error(w, fmt.Errorf("pageId missing or wrong %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	image, err := app.Store.CreateImage(imagePath, imagePath, uint(id))
// 	if err != nil || image == nil {
// 		http.Error(w, fmt.Errorf("error creating image %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	json, err := json.Marshal(image)

// 	if err != nil {
// 		http.Error(w, fmt.Errorf("could not encode output %s", err).Error(), http.StatusBadRequest)
// 		return
// 	}

// 	fmt.Fprintf(w, string(json))

// }
