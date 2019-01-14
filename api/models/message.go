package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

//Message model definition
type Message struct {
	ID      uint         `gorm:"primarykey"`
	Text    string       `valid:"alphanum,required"`
	Owner   Conversation `gorm:"foreignkey:OwnerID" valid:"-"`
	OwnerID uint         `valid:"numeric,required"`
}

//Validate model
func (m *Message) Validate(db *gorm.DB) {
	if m.Text == "" {
		db.AddError(errors.New("Text is empty"))
		return
	}

	if len(m.Text) > 255 {
		db.AddError(errors.New("invalid Text"))
		return
	}

	if m.OwnerID < 1 {
		db.AddError(errors.New("invalid OwnerID"))
		return
	}
}
