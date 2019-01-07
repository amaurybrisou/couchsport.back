package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	log "github.com/sirupsen/logrus"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
)

type FileStore struct {
	UploadPath, ImageBasePath string
}

func (app FileStore) DownloadImages(images []models.Image, ownerID uint) ([]models.Image, error) {
	if len(images) > 0 {
		for idx, i := range images {
			if i.File != "" && len(images) < 9 {
				filename, err := app.CreateFromB64(ownerID, i.URL, i.File)
				if err != nil {
					log.Error(err)
					return []models.Image{}, err
				}

				i.File = ""
				i.URL = filename
				i.Alt = filename
				images[idx] = i
			}
		}
	}
	return images, nil
}

func (app FileStore) CreateFromB64(UserID uint, b64, filename string) (string, error) {
	image, err := utils.B64ImageToFile(b64)
	if err != nil {
		log.Error(err)
		return "", err
	}

	path := strconv.FormatUint(uint64(UserID), 10) + "/"
	if _, err := os.Stat(app.UploadPath + path); os.IsNotExist(err) {
		log.Printf("creating directory %s", app.UploadPath+path)
		os.Mkdir(app.UploadPath+path, 0700)
	}

	if err != nil {
		log.Error(err)
		return "", err
	}

	//this is path which  we want to store the file
	f, err := os.OpenFile(app.UploadPath+path+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer f.Close()

	png.Encode(f, image)

	return app.ImageBasePath + path + filename, nil
}

func (app FileStore) FileUpload(r *http.Request, UserID uint) (string, error) {
	//this function returns the filename(to save in database) of the saved file or an error if it occurs
	r.ParseMultipartForm(32 << 20)

	//ParseMultipartForm parses a request body as multipart/form-data

	var fileName string

	file, handler, err := r.FormFile("file") //retrieve the file from form data
	defer file.Close()                       //close the file when we finish

	if err != nil {
		log.Error(err)
		return "", err
	}

	path := strconv.FormatUint(uint64(UserID), 10) + "/"
	if _, err := os.Stat(app.UploadPath + path); os.IsNotExist(err) {
		log.Printf("creating directory %s", app.UploadPath+path)
		os.Mkdir(app.UploadPath+path, 0700)
	}

	if err != nil {
		log.Error(err)
		return "", err
	}

	//this is path which  we want to store the file
	f, err := os.OpenFile(app.UploadPath+path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error(err)
		return "", err
	}

	fileName = handler.Filename
	defer f.Close()
	io.Copy(f, file)
	//here we save our file to our path

	return app.ImageBasePath + path + fileName, nil

}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
