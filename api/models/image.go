package models

import (
	"github.com/jinzhu/gorm"
)

//Image model definition
type Image struct {
	gorm.Model
	URL     string
	Alt     string
	File    string `gorm:"-"`
	OwnerID uint
}
