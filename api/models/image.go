package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

//Image model definition
type Image struct {
	gorm.Model
	URL     string `valid:"requri,required" gorm:"unique_index"`
	Alt     string `valid:"text"`
	File    string `gorm:"-"`
	OwnerID uint
}

//IsValid tells is Image is valid
func (image *Image) IsValid() bool {
	if image.URL == "" {
		return false
	}

	if len(image.URL) > 255 && image.File == "" {
		return false
	}

	if len(image.URL) > 255 && image.ID > 0 {
		return false
	}

	if len(image.Alt) > 255 {
		return false
	}

	return true
}

//Validate tells is Image is valid
func (image *Image) Validate(db *gorm.DB) {
	if image.URL == "" {
		db.AddError(errors.New("URL is empty"))
		return
	}

	if len(image.URL) > 255 {
		db.AddError(errors.New("invalid URL"))
		return
	}

	if len(image.Alt) > 255 {
		db.AddError(errors.New("invalid Alt"))
		return
	}
}
