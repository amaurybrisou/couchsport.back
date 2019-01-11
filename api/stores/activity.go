package stores

import (
	"couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type activityStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (cs activityStore) Migrate() {
	cs.Db.AutoMigrate(&models.Activity{})

	activities := []string{"Alpinisme", "Apnée", "Acrosport", "Badminton", "Base jump", "BMX", "Skateboard", "Escalade", "Canoë-kayak", "Canyoning", "Course", "Course d'orientation", "Crosse", "Cyclisme", "Danse", "Équitation", "Football", "Surf", "Golf", "Handball", "Kitesurfing", "Marathon", "Paddle", "Pêche", "Rafting", "Roller", "Ski Alpin", "Ski de fond", "Ski nordique", "Ski nautique", "Snowboard", "Tennis", "Tir à l'arc", "ULM", "Wakeboard", "Yoga"}

	for _, a := range activities {
		cs.Db.FirstOrCreate(&models.Activity{Name: a}, models.Activity{Name: a})
	}
}

//All Returns all the activities
func (cs activityStore) All() ([]models.Activity, error) {
	var activities []models.Activity
	if err := cs.Db.Find(&activities).Error; err != nil {
		return []models.Activity{}, err
	}
	return activities, nil
}
