package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

//Page model definition
type Page struct {
	gorm.Model
	Name, Description string  `valid:"text"`
	LongDescription   string  `gorm:"size:512;" valid:"text"`
	Images            []Image `gorm:"save_associations:true;foreignkey:OwnerID"`
	Lat               float64 `valid:"latitude"`
	Lng               float64 `valid:"longitude"`
	CouchNumber       *int    `valid:"numeric"`
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

//Validate  tells wheter the page is valid
func (page *Page) Validate(db *gorm.DB) {
	if !page.New && page.ID < 1 {
		db.AddError(errors.New("invalid PageID"))
		return
	}

	if page.Name == "" {
		db.AddError(errors.New("Name is empty"))
		return
	}

	if page.Description == "" {
		db.AddError(errors.New("Description is empty"))
		return
	}

	if (*page.CouchNumber) < 0 {
		db.AddError(errors.New("invalid CouchNumber"))
		return
	}

	if len(page.LongDescription) > 512 {
		db.AddError(errors.New("invalid LongDescription"))
		return
	}
}
