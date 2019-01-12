package stores

import (
	"couchsport/api/models"
	"fmt"
	"github.com/jinzhu/gorm"
)

type authorizationStore struct {
	Db *gorm.DB
}

//OwnImage tells you wheter the userID owns the imageID
func (me authorizationStore) OwnImage(userID, imageID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	profileID, err := me.GetProfileID(userID)
	if err != nil {
		return false, err
	}

	var count uint
	if err := me.Db.Model(models.Image{}).Where("owner_id = ?", profileID).Where("id = ?", imageID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v doesn't own this image %v", userID, imageID)
	}

	return true, nil
}

//OwnPage tells you wheter the userID owns the pageID
func (me authorizationStore) OwnPage(userID, pageID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	profileID, err := me.GetProfileID(userID)
	if err != nil {
		return false, err
	}

	var count uint
	if err := me.Db.Model(models.Page{}).Where("owner_id = ?", profileID).Where("id = ?", pageID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v doesn't own this page %v", userID, pageID)
	}

	return true, nil
}

//OwnProfile tells you wheter the userID owns the profile
func (me authorizationStore) OwnProfile(userID, profileID uint) (bool, error) {
	if profileID < 1 || userID < 1 {
		return false, fmt.Errorf("profileID or userID cannot be below 1")
	}

	profileID, err := me.GetProfileID(userID)
	if err != nil {
		return false, err
	}

	var count uint
	if err := me.Db.Model(models.Profile{}).Where("owner_id = ?", userID).Where("id = ?", profileID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v doesn't own this page %v", profileID, profileID)
	}

	return true, nil
}

//GetProfileID returns the profileID of the submitted userID
func (me authorizationStore) GetProfileID(userID uint) (uint, error) {
	var profile models.Profile
	if err := me.Db.Model(models.Profile{}).Where("owner_id = ?", userID).First(&profile).Error; err != nil {
		return 0, err
	}

	if profile.ID < 1 {
		return 0, fmt.Errorf("profileID cannot be below 1")
	}

	return profile.ID, nil
}
