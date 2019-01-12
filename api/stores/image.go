package stores

import (
	"couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type imageStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (app imageStore) Migrate() {
	app.Db.AutoMigrate(&models.Image{})
}

//All returns all the images in db
func (app imageStore) All() ([]models.Image, error) {
	var images []models.Image
	if err := app.Db.Find(&images).Error; err != nil {
		return []models.Image{}, err
	}
	return images, nil
}

//Delete an image by ID
func (app imageStore) Delete(imageID uint) (bool, error) {
	if err := app.Db.Exec("UPDATE images SET deleted_at = NOW() WHERE id = ? ", imageID).Error; err != nil {
		return false, err
	}

	return true, nil
}
