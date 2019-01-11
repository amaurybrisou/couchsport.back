package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
)

type profileStore struct {
	Db        *gorm.DB
	FileStore fileStore
}

//Migrate creates the table in database
func (me profileStore) Migrate() {
	me.Db.AutoMigrate(&models.Profile{})
}

//GetProfiles returns all profiles in database
func (me profileStore) All() ([]models.Profile, error) {
	var profiles []models.Profile
	if err := me.Db.Find(&profiles).Error; err != nil {
		return []models.Profile{}, err
	}
	return profiles, nil
}

//GetProfileByOwnerID returns the user profile
func (me profileStore) GetProfileByOwnerID(userID uint) (models.Profile, error) {
	var out = models.Profile{}
	if err := me.Db.Model(&models.Profile{}).Preload("Languages").Preload("OwnedPages").Preload("Activities").Where("owner_id = ?", userID).First(&out).Error; gorm.IsRecordNotFoundError(err) {
		out.OwnerID = userID
		if err := me.Db.Create(&out).Error; err != nil {
			return models.Profile{}, err
		}
	}
	return out, nil
}

//Update the profile
func (me profileStore) Update(profileID uint, profile models.Profile) (models.Profile, error) {
	if !profile.IsValid() {
		return models.Profile{}, fmt.Errorf("invalid profile")
	}

	if profile.AvatarFile != "" {
		filename, err := me.saveAvatar(profileID, profile.AvatarFile, profile.Avatar)
		if err != nil {
			return models.Profile{}, err
		}
		profile.AvatarFile = ""
		profile.Avatar = filename
	}

	if err := me.Db.Exec("DELETE FROM profile_languages WHERE profile_id = ?", profile.ID).Error; err != nil {
		return models.Profile{}, err
	}

	if err := me.Db.Exec("DELETE FROM profile_activities WHERE profile_id = ?", profile.ID).Error; err != nil {
		return models.Profile{}, err
	}

	if err := me.Db.Model(&models.Profile{}).Update(&profile).Error; err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (me profileStore) saveAvatar(profileID uint, filename, b64 string) (string, error) {
	//decode b64 string to bytes
	mime, buf, err := utils.B64ToImage(b64)
	if err != nil {
		return "", err
	}

	img, err := utils.ImageToTypedImage(mime, buf)
	if err != nil {
		return "", err
	}

	directory := "user-" + strconv.FormatUint(uint64(profileID), 10)
	filename, err = me.FileStore.Save(directory, filename, img)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (me profileStore) parseBody(tmp interface{}) (models.Profile, error) {
	r, ok := tmp.(models.Profile)

	if !ok {
		return models.Profile{}, fmt.Errorf("body is not of type Profile")
	}

	return r, nil
}
