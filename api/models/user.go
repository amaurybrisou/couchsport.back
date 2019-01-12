package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//User model definition
type User struct {
	gorm.Model
	Email     string `gorm:"unique_index"`
	Password  string
	Profile   Profile `gorm:"association_foreignkey:OwnerID"`
	ProfileID uint
	// FollowingPages  []*Page `gorm:"many2many:user_page_follower;"`
	// Friends         []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id;"`
	Type string
	New  bool `gorm:"-"`
}

//IsValid say wheter the underlying structure is a valid User
func (user *User) IsValid() bool {
	if user.Email == "" {
		return false
	}

	if user.Password == "" {
		return false
	}

	return true
}

//BeforeCreate generate the User ID, set Type to USER and hash the password
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Type", "USER")
	scope.SetColumn("Password", hashAndSalt([]byte(user.Password)))
	return nil
}

//AfterCreate empty the password column for security reasons, sets New to true and update Type to ADMIN if ID = 1
func (user *User) AfterCreate(scope *gorm.Scope) error {
	scope.SetColumn("Password", "")
	scope.SetColumn("New", true)
	if user.ID == 1 {
		scope.DB().Model(user).Update("type", "ADMIN")
	}
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
