package stores

import (
	"couchsport/api/models"
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type UserStore struct {
	Db *gorm.DB
}

func (app UserStore) Migrate() {
	app.Db.AutoMigrate(&models.User{})
}

func (app UserStore) GetUsers(keys url.Values) []models.User {
	var req = app.Db
	for i, v := range keys {
		switch i {
		case "profile":
			req = req.Preload("Profile")
		case "pages":
			req = req.Preload("OwnedPages")
		case "follow":
			req = req.Preload("FollowingPages")
		case "friends":
			req = req.Preload("Friends")
		case "id":
			req = req.Where("ID= ?", v)
		case "username":
			req = req.Where("username  LIKE ?", v)
		case "email":
			req = req.Where("email LIKE ?", v)
		}
	}

	var users []models.User
	if errs := req.Find(&users).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return users
}

func (app UserStore) FindOrCreate(u *models.User) (*models.User, error) {
	if ok, message := u.IsValid(); !ok {
		return nil, fmt.Errorf("Invalid credentials : %s", message)
	}

	if u.ID < 1 {
		//Page not found
		if err := app.Db.Create(u).Error; err != nil {
			log.Error(err)
			return nil, err
		}
		return u, nil
	}

	if err := app.Db.First(u).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return nil, nil
}

func (app UserStore) GetUser(formUser models.User) (models.User, error) {
	var outUser = models.User{}
	if errs := app.Db.Where(formUser).First(&outUser).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return outUser, nil
}

func (app UserStore) GetUserByID(userID uint) (models.User, error) {
	var outUser = models.User{}
	if errs := app.Db.Preload("Profile").Where("id = ?", userID).First(&outUser).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return outUser, nil
}
