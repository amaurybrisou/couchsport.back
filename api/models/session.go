package models

import (
	"time"
)

type Session struct {
	User      User
	UserID    uint
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
