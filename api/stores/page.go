package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io"
	"net/url"
)

type PageStore struct {
	Db           *gorm.DB
	FileStore    FileStore
	ImageStore   ImageStore
	ProfileStore ProfileStore
}

func (app PageStore) Migrate() {
	app.Db.AutoMigrate(&models.Page{})
}

func (app PageStore) GetPages(keys url.Values) []models.Page {
	var req = app.Db

	req = req.Preload("Images").Preload("Activities")

	for i, v := range keys {
		switch i {
		case "followers":
			req = req.Preload("Followers")
		case "profile":
			req = req.Preload("Owner")
		case "id":
			req = req.Where("ID= ?", v)
		}
	}

	var pages []models.Page
	if err := req.Find(&pages).Error; err != nil {
		log.Error(err)
	}
	return pages
}

func (app PageStore) CreateOrUpdate(userID uint, body io.Reader) (*models.Page, error) {
	pageObj, err := app.parseBody(body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	profileID, err := app.getProfileID(userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	pageObj.OwnerID = profileID
	savedPageObj := *pageObj

	if err := app.Db.Preload("Owner").Where("owner_id = ?", profileID).First(pageObj).Error; gorm.IsRecordNotFoundError(err) {
		//Page not found
		if err := app.Db.Create(pageObj).Error; err != nil {
			log.Error(err)
			return nil, err
		}

	}

	log.Println(pageObj.New)
	if !pageObj.New {
		pageObj = &savedPageObj
	}

	images, err := app.FileStore.DownloadImages(profileID, "page-", (*pageObj).Images)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err := app.Db.Exec("DELETE FROM page_activities WHERE page_id = ?", pageObj.ID).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	pageObj.Images = images // set parsed images from frontent
	if err := app.Db.Set("gorm:save_associations", true).Model(&models.Page{}).Update(pageObj).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return pageObj, nil

}

func (app PageStore) DeleteImage(userID uint, body io.Reader) (bool, error) {
	profileID, err := app.getProfileID(userID)
	if err != nil {
		log.Error(err)
		return false, err
	}

	result, err := app.ImageStore.Delete(profileID, body)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return result, nil
}

func (app PageStore) Delete(body io.Reader) (bool, error) {
	pageObj, err := app.parseBody(body)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if err := app.Db.Delete(pageObj).Error; err != nil {
		log.Error(err)
		return false, err
	}

	return true, nil
}

func (app PageStore) Publish(userID uint, body io.Reader) (bool, error) {
	pageObj, err := app.parseBody(body)
	if err != nil {
		log.Error(err)
		return false, err
	}

	profileID, err := app.getProfileID(userID)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if err := app.Db.Model(&models.Page{}).Where("id = ?", pageObj.ID).Where("owner_id = ? ", profileID).Update("Public", pageObj.Public).Error; err != nil {
		log.Error(err)
		return false, err

	}

	return true, nil
}

func (app PageStore) parseBody(body io.Reader) (*models.Page, error) {
	tmp, err := utils.ParseBody(&models.Page{}, body)
	if err != nil {
		return &models.Page{}, err
	}

	r, ok := tmp.(*models.Page)

	if !ok {
		return &models.Page{}, err
	}

	return r, nil
}

func (app PageStore) getProfileID(userID uint) (uint, error) {
	profile, err := app.ProfileStore.GetProfileByOwnerID(userID)
	if err != nil {
		log.Error(err)
		return uint(0), nil
	}

	return profile.ID, nil
}
