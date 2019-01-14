package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

//Conversation model definition
type Conversation struct {
	gorm.Model
	From     User `gorm:"foreignkey:FromID"`
	FromID   uint
	To       User `gorm:"foreignkey:ToID"`
	ToID     uint
	Messages []Message
}

//Validate model
func (c *Conversation) Validate(db *gorm.DB) {
	if c.FromID < 1 {
		db.AddError(errors.New("invalid c.FromID"))
	}

	if c.ToID < 1 {
		db.AddError(errors.New("invalid c.ToID"))
	}

	if len(c.Messages) < 1 {
		db.AddError(errors.New("invalid Messages"))
	}
}
