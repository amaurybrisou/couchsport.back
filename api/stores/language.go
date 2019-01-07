package stores

import (
	"couchsport/api/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type LanguageStore struct {
	Db *gorm.DB
}

func (cs LanguageStore) Migrate() {
	cs.Db.AutoMigrate(&models.Language{})
}

func (cs LanguageStore) GetLanguages() []models.Language {
	var languages []models.Language
	if errs := cs.Db.Find(&languages).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return languages
}
