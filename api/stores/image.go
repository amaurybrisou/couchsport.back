package stores

import (
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type imageStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (me imageStore) Migrate() {
	me.Db.AutoMigrate(&models.Image{})
	me.Db.Model(&models.Image{}).AddForeignKey("owner_id", "pages(id)", "CASCADE", "RESTRICT")
}

//All returns all the images in db
func (me imageStore) All() ([]models.Image, error) {
	var images []models.Image
	if err := me.Db.Find(&images).Error; err != nil {
		return []models.Image{}, err
	}
	return images, nil
}

//Delete an image by ID
func (me imageStore) Delete(imageID uint) (bool, error) {
	if err := me.Db.Exec("UPDATE images SET deleted_at = NOW() WHERE id = ? ", imageID).Error; err != nil {
		return false, err
	}

	return true, nil
}
