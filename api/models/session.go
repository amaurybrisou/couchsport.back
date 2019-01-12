package models

import (
	"time"
)

//Session model definition
type Session struct {
	Owner     User `gorm:"foreign_key:OwnerId;association_autoupdate:false;association_autocreate:false"`
	OwnerID   uint
	SessionID string
	Expires   time.Time `gorm:"default=now"`
	Validity  uint
}

//HasExpired determines wether the current session has expired or not
func (s Session) HasExpired() bool {
	now := time.Now()

	if s.Expires.After(now) {
		return false
	}
	return true
}
