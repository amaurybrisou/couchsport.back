package stores

import (
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/goland-amaurybrisou/couchsport/api/utils"
	"github.com/jinzhu/gorm"

	"net/url"
)

type userStore struct {
	Db *gorm.DB
}

func (me userStore) Migrate() {
	me.Db.AutoMigrate(&models.User{})
	me.Db.Model(&models.User{}).AddForeignKey("profile_id", "profiles(id)", "CASCADE", "CASCADE")
}

func (me userStore) All(keys url.Values) ([]models.User, error) {
	var req = me.Db
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

func (me userStore) New(user models.User) (models.User, error) {
	user.New = true

	var count int
	if err := me.Db.Model(&user).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		return models.User{}, err
	}

	if count > 0 {
		return models.User{}, fmt.Errorf("user already exist")
	}

	if err := me.Db.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

//GetProfile returns the user profile
func (me userStore) GetProfile(userID uint) (models.Profile, error) {
	var out = models.User{}
	if err := me.Db.Preload("Profile").Preload("Profile.Languages").Preload("Profile.Activities").Where("id = ?", userID).First(&out).Error; err != nil { //gorm.IsRecordNotFoundError(err) {
		return out.Profile, err
	}
	return out.Profile, nil
}

func (me userStore) GetByID(userID uint) (models.User, error) {
	var outUser = models.User{}
	if err := me.Db.Model(&models.User{}).Preload("Profile").Where("id = ?", userID).First(&outUser).Error; err != nil {
		return models.User{}, err
	}
	return outUser, nil
}

func (me userStore) GetByEmail(email string, create bool) (models.User, error) {
	var outUser = models.User{}
	if err := me.Db.Where("email = ?", email).First(&outUser).Error; create && gorm.IsRecordNotFoundError(err) {
		return me.NewWithoutPassword(email)
	} else if err != nil {
		return models.User{}, err
	}
	return outUser, nil
}

//OwnImage tells you wheter the userID owns the imageID
func (me userStore) OwnImage(userID, imageID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	user, err := me.GetByID(userID)
	if err != nil {
		return false, err
	}

	var count uint
	if err := me.Db.Model(models.Image{}).Where("owner_id = ?", user.ProfileID).Where("id = ?", imageID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v doesn't own this image %v", userID, imageID)
	}

	return true, nil
}

//OwnPage tells you wheter the userID owns the pageID
func (me userStore) OwnPage(userID, pageID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	user, err := me.GetByID(userID)
	if err != nil {
		return false, err
	}

	var count uint
	if err := me.Db.Model(models.Page{}).Where("owner_id = ?", user.ProfileID).Where("id = ?", pageID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v doesn't own this page %v", userID, pageID)
	}

	return true, nil
}

//OwnConversation tells you wheter the userID owns the conversation
func (me userStore) OwnConversation(userID, conversationID uint) (bool, error) {
	if userID < 1 {
		return false, fmt.Errorf("userID cannot be below 1")
	}

	profileID, err := me.GetProfileID(userID)
	if err != nil {
		return false, err
	}

	var count uint
	if err := me.Db.Model(models.Conversation{}).Where("from_id = ? OR to_id = ?", profileID, profileID).Where("id = ?", conversationID).Count(&count).Error; err != nil {
		return false, err
	}

	if count < 1 {
		return false, fmt.Errorf("user %v isn't part of this conversation %v", userID, conversationID)
	}

	return true, nil
}

//OwnProfile tells you wheter the userID owns the profile
func (me userStore) OwnProfile(userID, profileID uint) (bool, error) {
	if profileID < 1 || userID < 1 {
		return false, fmt.Errorf("profileID or userID cannot be below 1")
	}

	user, err := me.GetByID(userID)
	if err != nil {
		return false, err
	}

	if user.ID != profileID {
		return false, fmt.Errorf("user %v doesn't own this profile %v", profileID, profileID)
	}

	return true, nil
}

//GetProfileID returns the profileID of the submitted userID
func (me userStore) GetProfileID(userID uint) (uint, error) {
	profile, err := me.GetProfile(userID)
	return profile.ID, err
}

func parseBody(tmp interface{}) (models.User, error) {
	fmt.Println(tmp)
	r, ok := tmp.(*models.User)

	if !ok {
		return models.User{}, fmt.Errorf("body is not of type User")
	}

	return *r, nil
}

func (me userStore) NewWithoutPassword(email string) (models.User, error) {
	password := utils.RandStringBytesMaskImprSrc(len(email))
	user := models.User{
		Email:    email,
		Password: password,
	}

	return me.New(user)
}
