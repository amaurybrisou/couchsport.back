package server

import (
	"couchsport/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

func mustOpenDb(c *config.Config) *gorm.DB {
	log.Println(c.DriverName, c.DataSourceName)
	db, err := gorm.Open(c.DriverName, c.DataSourceName+"?"+c.DatabaseParams)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)

	return db
}
