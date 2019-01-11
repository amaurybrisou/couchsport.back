package handlers

import (
	"couchsport/api/models"
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type imageHandler struct {
	Store *stores.StoreFactory
}

//SoftDelete is called to set DeletedAt field to Now, not deleting the image
func (me imageHandler) SoftDelete(userID uint, w http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Body != nil {
		defer r.Body.Close()
	}

	image, err := me.parseBody(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	owns, err := me.Store.AuthorizationStore().OwnImage(userID, image.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	if !owns {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusForbidden)
		return
	}

	result, err := me.Store.ImageStore().Delete(image.ID)
	if err != nil {
		log.Error(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(struct{ Result bool }{Result: result})

	if err != nil {
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

func (me imageHandler) parseBody(body io.Reader) (models.Image, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return models.Image{}, err
	}

	var obj models.Image
	err = json.Unmarshal(b, &obj)

	if err != nil {
		return models.Image{}, err
	}

	return obj, nil
}
