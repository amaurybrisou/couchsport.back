package models

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	Username, Country, Gender, City, ZipCode, StreetName, Phone, Firstname, Lastname string `gorm:"type:varchar(50);"`
	Avatar                                                                           string
	AvatarFile                                                                       string `gorm:"-"`
	StreetNumber                                                                     uint
	UserID                                                                           uint
	Activities                                                                       []*Activity `gorm:"many2many:profile_activities;"`
	Languages                                                                        []*Language `gorm:"many2many:profile_languages;association_autoupdate:false;association_autocreate:false"`
}
