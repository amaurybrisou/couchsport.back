package stores

import (
	"couchsport/api/types"
	"couchsport/api/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
)

//FileStore manages the FileSystem
//For more details about PublicPath, ImageBasePath, FilePrefix please read the code
//FileSystem argument specify which filesystem you want to use
type FileStore struct {
	PublicPath, ImageBasePath, FilePrefix string
	FileSystem                            types.FileSystem
}

//Save a file on the filesystem at path computed from ImageBasePath + prefix + UserID
func (app FileStore) Save(UserID uint, prefix, filename string, buf io.Reader) (string, error) {

	if UserID < 1 {
		err := fmt.Errorf("userID is incorrect")
		log.Error(err)
		return "", err
	}

	if filename == "" {
		err := fmt.Errorf("filename is incorrect")
		log.Error(err)
		return "", err
	}

	path := app.ImageBasePath + prefix + strconv.FormatUint(uint64(UserID), 10) + "/"

	path, err := utils.CreateDirIfNotExists(app.FileSystem, app.PublicPath+path)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Printf("Openning file %s", app.PublicPath+path+app.FilePrefix+filename)
	f, err := app.FileSystem.OpenFile(path + app.FilePrefix + filename)
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

	log.Printf("%d bytes wrote at %s", count, path+filename)

	return "/" + path + app.FilePrefix + filename, nil
}
