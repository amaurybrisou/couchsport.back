package stores

import (
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type activityStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (me activityStore) Migrate() {
	me.Db.AutoMigrate(&models.Activity{})

	me.Db.Table("page_activities").AddForeignKey("activity_id", "activities(id)", "NO ACTION", "NO ACTION")
	me.Db.Table("page_activities").AddForeignKey("page_id", "pages(id)", "CASCADE", "NO ACTION")
	me.Db.Table("page_activities").AddUniqueIndex("activity_id_page_id_unique", "page_id, activity_id")

	me.Db.Table("profile_activities").AddForeignKey("activity_id", "activities(id)", "NO ACTION", "NO ACTION")
	me.Db.Table("profile_activities").AddForeignKey("profile_id", "profiles(id)", "CASCADE", "NO ACTION")
	me.Db.Table("profile_activities").AddUniqueIndex("activity_id_profile_id_unique", "profile_id, activity_id")

	activities := []string{"Alpinisme", "Apnée", "Acrosport", "Badminton", "Base jump", "BMX", "Skateboard", "Escalade", "Canoë-kayak", "Canyoning", "Course", "Course d'orientation", "Crosse", "Cyclisme", "Danse", "Équitation", "Football", "Surf", "Golf", "Handball", "Kitesurfing", "Marathon", "Paddle", "Pêche", "Rafting", "Roller", "Ski Alpin", "Ski de fond", "Ski nordique", "Ski nautique", "Snowboard", "Tennis", "Tir à l'arc", "ULM", "Wakeboard", "Yoga"}

	for _, a := range activities {
		me.Db.FirstOrCreate(&models.Activity{Name: a}, models.Activity{Name: a})
	}

}

//All Returns all the activities
func (me activityStore) All() ([]models.Activity, error) {
	var activities []models.Activity
	if err := me.Db.Find(&activities).Error; err != nil {
		return []models.Activity{}, err
	}
	return activities, nil
}
