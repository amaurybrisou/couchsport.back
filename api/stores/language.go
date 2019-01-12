package stores

import (
	"couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type languageStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (cs languageStore) Migrate() {
	cs.Db.AutoMigrate(&models.Language{})
}

//All returns all the languages in db
func (cs languageStore) All() ([]models.Language, error) {
	var languages []models.Language
	if err := cs.Db.Find(&languages).Error; err != nil {
		return []models.Language{}, err
	}
	return languages, nil
}
