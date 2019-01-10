package types

import (
	"io"
	"os"
)

type FileSystem interface {
	OpenFile(name string) (io.WriteCloser, error)
	MkdirAll(path string) error
	Stat(name string) (os.FileInfo, error)
	IsNotExist(error) bool
}
