package config

import (
	"log"
	"os"
)

type logger struct {
	Name, Mode, FilePath string
}

//Logger creates an initial logger base on specified configuration
//Name is the Name of the logger
//Mode helps this function to decide beween file/stderr/stdout buffer
//FilePath is the fileName output if "file" is specified in Mode
func Logger(l logger) *log.Logger {
	var output = os.Stdout
	var err error

	switch l.Mode {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "file":
		output, err = os.OpenFile(l.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panic(err)
		}
	}

	return log.New(output, l.Name, log.Lshortfile|log.Ldate|log.Ltime)
}
