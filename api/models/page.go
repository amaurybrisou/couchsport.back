package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

//Page model definition
type Page struct {
	gorm.Model
	Name, Description string
	LongDescription   string  `gorm:"size:512;"`
	Images            []Image `gorm:"save_associations:false;foreignkey:OwnerID"`
	Lat               float64
	Lng               float64
	Followers         []*User `gorm:"many2many:user_page_follower"`
	Owner             Profile `gorm:"foreignkey:OwnerID;association_autoupdate:false;association_autocreate:false"`
	OwnerID           uint
	Public            bool        `gorm:"default:1"`
	Activities        []*Activity `gorm:"many2many:page_activities;association_autoupdate:false;association_autocreate:false"`
	New               bool        `gorm:"-"`
}

//BeforeCreate is a gorm hook
func (page *Page) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

//AfterCreate is a gorm hook
func (page *Page) AfterCreate(scope *gorm.Scope) error {
	scope.SetColumn("New", true)
	return nil
}

//IsValid  tells wheter the page is valid
func (page *Page) IsValid(state string) bool {
	if state == "UPDATE" && page.ID < 1 {
		return false
	}

	if page.Name == "" {
		return false
	}
	if page.Description == "" {
		return false
	}

	return true
}
