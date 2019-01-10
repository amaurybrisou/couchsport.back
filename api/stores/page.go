package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"io"
	"net/url"
	"strconv"
)

type PageStore struct {
	Db           *gorm.DB
	FileStore    FileStore
	ImageStore   ImageStore
	ProfileStore ProfileStore
}

//Migrate creates the model schema in database
func (app PageStore) Migrate() {
	app.Db.AutoMigrate(&models.Page{})
}

//GetPages returns all pages in Database
//Additional keys (url.Values) can be specified :
//followers : returns pages followers
//profile : returns pages profiles
//id: fetch a specific page
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

func (app PageStore) GetPagesByOwnerID(userID uint) ([]models.Page, error) {

	profileID, err := app.getProfileID(userID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var pages []models.Page
	if err := app.Db.Model(models.Page{}).Preload("Activities").Preload("Images").Where("owner_id = ?", profileID).Find(&pages).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return pages, nil
}

//CreateOrUpdate creates or update a page owned by userID. body is models.Page JSON encoded as io.Reader
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

	pageObj.Images = []models.Image{}
	if pageObj.ID == 0 {
		if err := app.Db.Create(pageObj).Error; err != nil {
			log.Error(err)
			return nil, err
		}
	}

	if err := app.Db.Where("owner_id = ?", profileID).First(pageObj).Error; gorm.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("you are not the owner of the page, thus cannot edit page %d", pageObj.ID)
	}

	if !pageObj.New {
		pageObj.Name = savedPageObj.Name
		pageObj.Images = savedPageObj.Images
		pageObj.Description = savedPageObj.Description
		pageObj.LongDescription = savedPageObj.LongDescription
		pageObj.Lat = savedPageObj.Lat
		pageObj.Lng = savedPageObj.Lng
	}

	directory := "page-" + strconv.FormatUint(uint64(profileID), 10)
	images, err := app.downloadImages(directory, (*pageObj).Images)
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

//DeleteImage set image.DeletedAt field at time.Now(); soft delete thus
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

//Delete set page.DeletedAt to time.Now() // soft delete thus
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

//Publish set page.Public field to 0 or 1
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

func (app PageStore) downloadImages(directory string, images []models.Image) ([]models.Image, error) {
	if len(images) > 0 {
		for idx, i := range images {
			if i.File != "" && len(images) < 9 {

				//decode b64 string to bytes
				mime, buf, err := utils.B64ToImage(i.URL)
				if err != nil {
					log.Error(err)
					return []models.Image{}, err
				}

				img, err := utils.ImageToTypedImage(mime, buf)
				if err != nil {
					log.Error(err)
					return []models.Image{}, err
				}

				filename, err := app.FileStore.Save(directory, i.File, img)
				if err != nil {
					log.Error(err)
					return []models.Image{}, err
				}

				i.File = ""
				i.URL = filename
				images[idx] = i
			}
		}
	}
	return images, nil
}
