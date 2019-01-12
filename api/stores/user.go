package stores

import (
	"couchsport/api/models"
	"fmt"
	"github.com/jinzhu/gorm"

	"net/url"
)

type userStore struct {
	Db *gorm.DB
}

func (app userStore) Migrate() {
	app.Db.AutoMigrate(&models.User{})
}

func (app userStore) All(keys url.Values) ([]models.User, error) {
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
	if err := req.Find(&users).Error; err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (app userStore) New(user models.User) (models.User, error) {
	if !user.IsValid() {
		return models.User{}, fmt.Errorf("Invalid user")
	}

	var count int
	if err := app.Db.Model(models.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return models.User{}, err
	}

	if count > 0 {
		return models.User{}, fmt.Errorf("user already exist")
	}

	if err := app.Db.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (app userStore) GetByID(userID uint) (models.User, error) {
	var outUser = models.User{}
	if err := app.Db.Model(&models.User{}).Preload("Profile").Where("id = ?", userID).First(&outUser).Error; err != nil {
		return models.User{}, err
	}
	return outUser, nil
}

func (app userStore) GetByEmail(email string) (models.User, error) {
	var outUser = models.User{}
	if err := app.Db.Where("email = ?", email).First(&outUser).Error; err != nil {
		return models.User{}, err
	}
	return outUser, nil
}

func parseBody(tmp interface{}) (models.User, error) {
	fmt.Println(tmp)
	r, ok := tmp.(*models.User)

	if !ok {
		return models.User{}, fmt.Errorf("body is not of type User")
	}

	return *r, nil
}
