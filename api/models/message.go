package models

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"time"
)

//Message model definition
type Message struct {
	ID      uint      `gorm:"primarykey"`
	Date    time.Time `gorm:"default:NOW()"`
	Text    string    `valid:"text,required"`
	From    Profile   `gorm:"foreign_key:ProfileID"`
	FromID  uint      `gorm:"required"`
	OwnerID uint      `valid:"numeric,required"`
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

	if m.FromID < 1 {
		db.AddError(errors.New("invalid FromID"))
		return
	}
}

//SendMessageBodyModel model definition (model used when decoding body for in conversationHandler.HandleMessage)
type SendMessageBodyModel struct {
	Email string `valid:"email,required"`
	Text  string `valid:"text,required"`
	ToID  uint   `valid:"numeric"`
}

//Validate use govalidators to check expression values
func (me SendMessageBodyModel) Validate() (bool, error) {
	if me.ToID < 1 {
		return false, fmt.Errorf("invalid ToID: %v", me.ToID)
	}
	return govalidator.ValidateStruct(me)
}
