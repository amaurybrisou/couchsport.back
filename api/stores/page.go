package stores

import (
	"couchsport/api/models"
	"couchsport/api/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
	"strconv"
)

type pageStore struct {
	Db           *gorm.DB
	FileStore    fileStore
	ImageStore   imageStore
	ProfileStore profileStore
}

//Migrate creates the model schema in database
func (me pageStore) Migrate() {
	me.Db.AutoMigrate(&models.Page{})
}

//All returns all pages in Database
//Additional keys (url.Values) can be specified :
//followers : returns pages followers
//profile : returns pages profiles
//id: fetch a specific page
func (me pageStore) All(keys url.Values) ([]models.Page, error) {
	var req = me.Db

	req = req.Preload("Images").Preload("Activities")

	for i, v := range keys {
		switch i {
		case "followers":
			req = req.Preload("Followers")
		case "profile":
			req = req.Preload("Owner")
		case "id":
			req = req.Where("ID= ?", v)
		case "owner_id":
			req = req.Where("owner_id = ?", v)
		}
	}

	var pages []models.Page
	if err := req.Find(&pages).Error; err != nil {
		return []models.Page{}, err
	}
	return pages, nil
}

//GetPagesByOwnerID return all profile details
func (me pageStore) GetPagesByOwnerID(profileID uint) ([]models.Page, error) {
	var pages []models.Page
	if err := me.Db.Model(&models.Page{}).Preload("Activities").Preload("Images").Where("owner_id = ?", profileID).Find(&pages).Error; err != nil {
		return nil, err
	}

	return pages, nil
}

//New creates a page
func (me pageStore) New(profileID uint, page models.Page) (models.Page, error) {
	if !page.IsValid("NEW") {
		return models.Page{}, fmt.Errorf("invalid page")
	}

	page.OwnerID = profileID

	directory := "page-" + strconv.FormatUint(uint64(profileID), 10)
	images, err := me.downloadImages(directory, page.Images)
	if err != nil {
		return models.Page{}, err
	}

	page.Images = images

	if err := me.Db.Create(&page).Error; err != nil {
		return models.Page{}, err
	}

	return page, nil
}

//Update the page
func (me pageStore) Update(userID uint, page models.Page) (models.Page, error) {
	if !page.IsValid("UPDATE") {
		return models.Page{}, fmt.Errorf("pageId cannot be below 0 %v", page.ID)
	}

	if len(page.Images) > 0 {
		directory := "page-" + strconv.FormatUint(uint64(userID), 10)
		images, err := me.downloadImages(directory, page.Images)
		if err != nil {
			return models.Page{}, err
		}
		me.Db.Model(&page).Association("Images").Replace(images) // update with newly parsed images and previous ones
	}

	me.Db.Unscoped().Table("page_activities").Where("activity_id NOT IN (?)", me.getActivitiesIDS(page.Activities)).Where("page_id = ?", page.ID).Delete(&models.Image{})
	me.Db.Model(&page).Association("Activities").Append(page.Activities)

	if err := me.Db.Set("gorm:save_associations", true).Model(&page).Update(&page).Error; err != nil {
		return models.Page{}, err
	}

	return page, nil
}

//Delete set page.DeletedAt to time.Now() // soft delete thus
func (me pageStore) Delete(userID, pageID uint) (bool, error) {
	if err := me.Db.Exec("DELETE FROM pages WHERE id = ?", pageID).Error; err != nil {
		return false, err
	}

	return true, nil
}

//Publish set page.Public field to 0 or 1
func (me pageStore) Publish(userID, pageID uint, status bool) (bool, error) {
	if err := me.Db.Model(&models.Page{}).Where("id = ?", pageID).Update("Public", status).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (me pageStore) getImagesIDS(images []models.Image) []uint {
	tmp := []uint{0}
	for _, el := range images {
		tmp = append(tmp, el.ID)
	}
	return tmp
}

func (me pageStore) getActivitiesIDS(activities []*models.Activity) []uint {
	tmp := []uint{0}
	for _, l := range activities {
		tmp = append(tmp, (*l).ID)
	}
	return tmp
}

func (me pageStore) downloadImages(directory string, images []models.Image) ([]models.Image, error) {
	if len(images) > 0 {
		for idx, i := range images {
			if i.File != "" && len(images) < 9 {

				//decode b64 string to bytes
				mime, buf, err := utils.B64ToImage(i.URL)
				if err != nil {
					return []models.Image{}, err
				}

				img, err := utils.ImageToTypedImage(mime, buf)
				if err != nil {
					return []models.Image{}, err
				}

				filename, err := me.FileStore.Save(directory, i.File, img)
				if err != nil {
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
