package server

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amaurybrisou/couchsport.back/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func mustOpenDb(c *config.Config) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	log.Println(c.DriverName, c.DataSourceName)
	dsn := fmt.Sprintf("%s?%s", c.DataSourceName, c.DatabaseParams)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		log.Fatal(err)
	}

	// db.LogMode(c.Verbose)

	// validations.RegisterCallbacks(db)

	return db
}
