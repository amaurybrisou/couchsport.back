package models

type Activity struct {
	ID       uint       `gorm:"primarykey"`
	Name     string     `gorm:"unique_index"`
	Profiles []*Profile `gorm:"many2many:profile_activities;"`
	Pages    []*Page    `gorm:"many2many:page_activities;"`
}
