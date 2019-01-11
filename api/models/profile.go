package models

import "github.com/jinzhu/gorm"

//Profile definition
type Profile struct {
	gorm.Model
	Username, Country, Gender, City, ZipCode, StreetName, Phone, Firstname, Lastname string `gorm:"type:varchar(50);"`
	Avatar                                                                           string
	AvatarFile                                                                       string `gorm:"-"`
	StreetNumber                                                                     uint
	OwnerID                                                                          uint        `gorm:"association_autoupdate:false;association_autocreate:false"`
	OwnedPages                                                                       []Page      `gorm:"foreignkey:OwnerID;association_autoupdate:false;association_autocreate:false"`
	Activities                                                                       []*Activity `gorm:"many2many:profile_activities;association_autoupdate:false;association_autocreate:false"`
	Languages                                                                        []*Language `gorm:"many2many:profile_languages;association_autoupdate:false;association_autocreate:false"`
}

//IsValid telle wether the profile is valid
func (p *Profile) IsValid() bool {
	if p.ID < 1 {
		return false
	}

	return true
}
