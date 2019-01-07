package models

import (
	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	URL    string
	Alt    string
	File   string `gorm:"-"`
	PageID uint
}
