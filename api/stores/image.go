package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io"
)

type ImageStore struct {
	Db *gorm.DB
}

func (app ImageStore) Migrate() {
	app.Db.AutoMigrate(&models.Image{})
}

func (app ImageStore) GetImages() []models.Image {
	var images []models.Image
	if err := app.Db.Find(&images).Error; err != nil {
		log.Error(err)
	}
	return images
}

func (app ImageStore) Delete(profileID uint, body io.Reader) (bool, error) {
	image, err := app.parseBody(body)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if err := app.Db.Exec("UPDATE images AS i INNER JOIN pages AS p ON i.page_id = p.id SET i.deleted_at = NOW() WHERE i.id = ? AND p.owner_id = ?", image.ID, profileID).Error; err != nil {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (app ImageStore) parseBody(body io.Reader) (*models.Image, error) {
	tmp, err := utils.ParseBody(&models.Image{}, body)
	if err != nil {
		return &models.Image{}, err
	}

	r, ok := tmp.(*models.Image)

	if !ok {
		return &models.Image{}, err
	}

	return r, nil
}
