package stores

import (
	"couchsport/api/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type InvitationStore struct {
	Db *gorm.DB
}

func (cs InvitationStore) Migrate() {
	cs.Db.AutoMigrate(&models.Invitation{})
}

func (cs InvitationStore) GetInvitations(keys url.Values) []models.Invitation {
	var req = cs.Db
	for i, v := range keys {
		switch i {
		case "profile":
			req = req.Preload("From")
			req = req.Preload("To")
		case "status":
			req = req.Where("Status = ?", v)
		case "id":
			req = req.Where("ID= ?", v)
		case "from":
			req = req.Where("from_id = ?", v)
		case "to":
			req = req.Where("to_id = ?", v)
		}
	}

	var invitations []models.Invitation
	if errs := req.Find(&invitations).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return invitations
}
