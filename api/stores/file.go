package stores

import (
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/types"
	"github.com/goland-amaurybrisou/couchsport/api/utils"
	log "github.com/sirupsen/logrus"
	"io"
)

type fileStore struct {
	PublicPath, ImageBasePath, FilePrefix string
	FileSystem                            types.FileSystem
}

//Save a file on the filesystem at path computed from ImageBasePath + directory + UserID
//directory is prepend before userId
func (me fileStore) Save(directory, filename string, buf io.Reader) (string, error) {

	if filename == "" {
		err := fmt.Errorf("filename is incorrect")
		return "", err
	}

	path := me.ImageBasePath
	if directory != "" {
		path += "/" + directory + "/"
	}

	fsPath, err := utils.CreateDirIfNotExists(me.FileSystem, me.PublicPath+"/"+path)
	if err != nil {
		return "", err
	}

	log.Printf("Openning file %s", fsPath+me.FilePrefix+filename)
	f, err := me.FileSystem.OpenFile(fsPath + me.FilePrefix + filename)
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

	log.Printf("%d bytes wrote at %s", count, fsPath+me.FilePrefix+filename)

	return "/" + path + me.FilePrefix + filename, nil
}
