package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io"
)

type ProfileStore struct {
	Db        *gorm.DB
	FileStore FileStore
}

func (app ProfileStore) Migrate() {
	app.Db.AutoMigrate(&models.Profile{})
}

func (app ProfileStore) GetProfiles() []models.Profile {
	var profiles []models.Profile
	if errs := app.Db.Find(&profiles).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return profiles
}

func (app ProfileStore) GetProfileByOwnerID(userID uint) (models.Profile, error) {
	var out = models.Profile{}
	if err := app.Db.Preload("Languages").Preload("Activities").Where("user_id = ?", userID).First(&out).Error; gorm.IsRecordNotFoundError(err) {
		log.Error(err)

		out.UserID = userID
		if err := app.Db.Save(&out).Error; err != nil {
			log.Error(err)
		}
	}
	return out, nil
}

func (app ProfileStore) Update(userID uint, body io.Reader) (*models.Profile, error) {
	profile, err := app.parseBody(body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if profile.AvatarFile != "" {
		filename, err := app.FileStore.CreateFromB64(userID, "user-", profile.AvatarFile, profile.Avatar)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		profile.AvatarFile = ""
		profile.Avatar = filename
	}

	if err := app.Db.Exec("DELETE FROM profile_languages WHERE profile_id = ?", profile.ID).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	if err := app.Db.Exec("DELETE FROM profile_activities WHERE profile_id = ?", profile.ID).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	if err := app.Db.Model(profile).Where("user_id = ?", userID).Update(profile).Error; err != nil {
		log.Error(err)
		return nil, err
	}
	return profile, nil
}

func (app ProfileStore) parseBody(body io.Reader) (*models.Profile, error) {
	tmp, err := utils.ParseBody(&models.Profile{}, body)
	if err != nil {
		return &models.Profile{}, err
	}

	r, ok := tmp.(*models.Profile)

	if !ok {
		return &models.Profile{}, err
	}

	return r, nil
}
