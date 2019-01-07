package models

import (
	"github.com/jinzhu/gorm"
)

type Invitation struct {
	gorm.Model
	From    User `gorm:"foreignkey:FromID"`
	To      User `gorm:"foreignkey:ToID"`
	Status  string
	Message string
	FromID  uint
	ToID    uint
}
