package stores

import (
	"couchsport/api/types"
	"couchsport/api/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

//FileStore manages the FileSystem
//For more details about PublicPath, ImageBasePath, FilePrefix please read the code
//FileSystem argument specify which filesystem you want to use
type FileStore struct {
	PublicPath, ImageBasePath, FilePrefix string
	FileSystem                            types.FileSystem
}

//Save a file on the filesystem at path computed from ImageBasePath + directory + UserID
//directory is prepend before userId
func (app FileStore) Save(directory, filename string, buf io.Reader) (string, error) {

	if filename == "" {
		err := fmt.Errorf("filename is incorrect")
		log.Error(err)
		return "", err
	}

	path := app.ImageBasePath
	if directory != "" {
		path += directory + "/"
	}

	fsPath, err := utils.CreateDirIfNotExists(app.FileSystem, app.PublicPath+path)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Printf("Openning file %s", fsPath+app.FilePrefix+filename)
	f, err := app.FileSystem.OpenFile(fsPath + app.FilePrefix + filename)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer f.Close()
	count, err := io.Copy(f, buf)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("file not created, no data to write")
	}

	log.Printf("%d bytes wrote at %s", count, fsPath+app.FilePrefix+filename)

	return "/" + path + app.FilePrefix + filename, nil
}
