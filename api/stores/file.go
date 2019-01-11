package stores

import (
	"couchsport/api/types"
	"couchsport/api/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

type fileStore struct {
	PublicPath, ImageBasePath, FilePrefix string
	FileSystem                            types.FileSystem
}

//Save a file on the filesystem at path computed from ImageBasePath + directory + UserID
//directory is prepend before userId
func (app fileStore) Save(directory, filename string, buf io.Reader) (string, error) {

	if filename == "" {
		err := fmt.Errorf("filename is incorrect")
		return "", err
	}

	path := app.ImageBasePath
	if directory != "" {
		path += directory + "/"
	}

	fsPath, err := utils.CreateDirIfNotExists(app.FileSystem, app.PublicPath+path)
	if err != nil {
		return "", err
	}

	log.Printf("Openning file %s", fsPath+app.FilePrefix+filename)
	f, err := app.FileSystem.OpenFile(fsPath + app.FilePrefix + filename)
	if err != nil {
		return "", err
	}

	defer f.Close()
	count, err := io.Copy(f, buf)
	if err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("file not created, no data to write")
	}

	log.Printf("%d bytes wrote at %s", count, fsPath+app.FilePrefix+filename)

	return "/" + path + app.FilePrefix + filename, nil
}
