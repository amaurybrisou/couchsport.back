package models

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
)

//Conversation model definition
type Conversation struct {
	gorm.Model
	From     Profile   `gorm:"foreignkey:FromID"`
	FromID   uint      `gorm:"association_autoupdate:false;;association_autosave:false;save_associations:false;association_save_reference:false"`
	To       Profile   `gorm:"foreignkey:ToID"`
	ToID     uint      `gorm:"association_autoupdate:false;;association_autosave:false;save_associations:false;association_save_reference:false"`
	Messages []Message `gorm:"foreignkey:OwnerID"`
	New      bool      `gorm:"-"`
}

//Validate model
func (me *Conversation) Validate(db *gorm.DB) {
	if me.FromID < 1 {
		db.AddError(errors.New("invalid FromID"))
	}

	if me.ToID < 1 {
		db.AddError(errors.New("invalid ToID"))
	}

	if me.ToID == me.FromID {
		db.AddError(errors.New("invalid Conversation"))
	}

	if len(me.Messages) < 1 && me.ID < 1 && !me.New {
		db.AddError(errors.New("invalid Messages"))
	}
}

//AddMessage to the expression Messages
func (me *Conversation) AddMessage(fromID, toID uint, text string) Message {
	m := Message{Text: text, OwnerID: me.ID, FromID: fromID, ToID: toID}
	me.Messages = append(me.Messages, m)
	return m
}

//ToJSON converts model to a json string
func (me *Conversation) ToJSON() (string, error) {
	j, err := json.Marshal(me)
	if err != nil {
		return "", err
	}
	return string(j), nil
}
