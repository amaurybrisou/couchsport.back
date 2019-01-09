package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Page struct {
	gorm.Model
	Name, Description string
	LongDescription   string  `gorm:"size:512;"`
	Images            []Image `gorm:"save_associations:false"`
	Lat               float64
	Lng               float64
	Followers         []*User `gorm:"many2many:user_page_follower"`
	Owner             Profile `gorm:"foreign_key:OwnerId"`
	OwnerID           uint
	Public            bool        `gorm:"default:1"`
	Activities        []*Activity `gorm:"many2many:page_activities;"`
	New               bool        `gorm:"-"`
}

func (page *Page) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}
func (page *Page) AfterCreate(scope *gorm.Scope) error {
	scope.SetColumn("New", true)
	return nil
}
