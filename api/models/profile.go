package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

//Profile definition
type Profile struct {
	gorm.Model
	Username, Country, Gender, City, StreetName, Firstname, Lastname string `valid:"name" gorm:"type:varchar(50);"`
	Phone                                                            string `valid:"alphanum"`
	ZipCode                                                          string `valid:"zipcode"`
	Avatar                                                           string `valid:"requri"`
	AvatarFile                                                       string `gorm:"-" valid:"-"`
	StreetNumber                                                     uint   `valid:"numeric"`
	// User                                                                             User
	// OwnerID                                                                          uint        `gorm:"association_autoupdate:false;association_autocreate:false"`
	OwnedPages    []Page         `gorm:"foreignkey:OwnerID;association_autoupdate:false;association_autocreate:false"`
	Conversations []Conversation `gorm:"association_autoupdate:false;association_autocreate:false"`

	Activities []*Activity `gorm:"many2many:profile_activities;association_autoupdate:false;association_autocreate:false"`
	Languages  []*Language `gorm:"many2many:profile_languages;association_autoupdate:false;association_autocreate:false"`
}

//Validate model
func (p *Profile) Validate(db *gorm.DB) {
	if len(p.Username) > 255 {
		db.AddError(errors.New("Username invalid"))
		return
	}

	if len(p.Firstname) > 255 {
		db.AddError(errors.New("Firstname invalid"))
		return
	}

	if len(p.Lastname) > 255 {
		db.AddError(errors.New("Lastname invalid"))
		return
	}

	if len(p.Phone) > 255 {
		db.AddError(errors.New("Phone invalid"))
		return
	}

	if len(p.ZipCode) > 255 {
		db.AddError(errors.New("ZipCode invalid"))
		return
	}

	if len(p.StreetName) > 255 {
		db.AddError(errors.New("StreetName invalid"))
		return
	}

	if len(p.Country) > 255 {
		db.AddError(errors.New("Country invalid"))
		return
	}

	if len(p.City) > 255 {
		db.AddError(errors.New("City invalid"))
		return
	}

	if len(p.Gender) > 1 {
		db.AddError(errors.New("Gender invalid"))
		return
	}

}
