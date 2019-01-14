package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

//User model definition
type User struct {
	gorm.Model
	Email     string  `gorm:"unique_index" valid:"email,required"`
	Password  string  `valid:"required,length(8|255)"`
	Profile   Profile `valid:"-" gorm:"foreignkey:ProfileID;association_autocreate:false;save_associations:false;association_save_reference:true;"`
	ProfileID uint    `valid:"numeric"`
	// // FollowingPages  []*Page `gorm:"many2many:user_page_follower;"`
	// Friends         []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id;"`
	Type string `valid:"in(ADMIN|USER)"`
	New  bool   `gorm:"-" valid:"-"`
}

//Validate model
func (user User) Validate(db *gorm.DB) {
	if user.Email == "" {
		db.AddError(errors.New("email is empty"))
		return
	}

	if user.Password == "" {
		db.AddError(errors.New("password is empty"))
		return
	}

	if user.ProfileID == 0 && !user.New {
		db.AddError(errors.New("profileID invalid"))
		return
	}

	return
}

//BeforeCreate generate the User ID, set Type to USER and hash the password
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	profile := Profile{}
	if err := scope.DB().Create(&profile).Error; err != nil {
		return err
	}

	scope.SetColumn("Profile", profile)
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
