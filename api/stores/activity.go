package stores

import (
	"couchsport/api/models"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type ActivityStore struct {
	Db *gorm.DB
}

func (cs ActivityStore) Migrate() {
	cs.Db.AutoMigrate(&models.Activity{})

	activities := []string{"Alpinisme", "Apnée", "Acrosport", "Badminton", "Base jump", "BMX", "Skateboard", "Escalade", "Canoë-kayak", "Canyoning", "Course", "Course d'orientation", "Crosse", "Cyclisme", "Danse", "Équitation", "Football", "Surf", "Golf", "Handball", "Kitesurfing", "Marathon", "Paddle", "Pêche", "Rafting", "Roller", "Ski Alpin", "Ski de fond", "Ski nordique", "Ski nautique", "Snowboard", "Tennis", "Tir à l'arc", "ULM", "Wakeboard", "Yoga"}

	for _, a := range activities {
		cs.Db.FirstOrCreate(&models.Activity{Name: a}, models.Activity{Name: a})
	}
}

func (cs ActivityStore) GetActivities() []models.Activity {
	var activities []models.Activity
	if errs := cs.Db.Find(&activities).GetErrors(); len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
	}
	return activities
}
