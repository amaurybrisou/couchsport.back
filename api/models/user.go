package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email, Password string
	// FollowingPages  []*Page `gorm:"many2many:user_page_follower;"`
	// Friends         []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id;"`
	Profile   Profile `gorm:"association_foreignkey:UserID"`
	ProfileID uint
	Type      string
	New       bool `gorm:"-"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.NewGen())
	scope.SetColumn("Type", "USER")
	scope.SetColumn("Password", hashAndSalt([]byte(user.Password)))
	return nil
}

func (user *User) AfterCreate(scope *gorm.Scope) error {
	scope.SetColumn("New", true)
	return nil
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
