package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Config describes the required configuration
type Config struct {
	Name                                                                string
	Listen                                                              string
	Port                                                                int
	Username, Password, DataFile, PublicPath, UploadPath, ImageBasePath string
	DataSourceName, DatabaseParams, DriverName, FixtureFile             string
	Logger                                                              struct {
		Name, Mode, FilePath string
	}
}

//Load loads the configuration according to env parameter. i.e config.dev.json
func Load(env string) *Config {
	jsonFile, err := ioutil.ReadFile("config." + env + ".json")
	if err != nil {
		log.Fatal(err)
	}

	var config *Config

	json.Unmarshal([]byte(jsonFile), &config)

	return config
}
