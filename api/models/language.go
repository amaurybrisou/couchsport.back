package models

type Language struct {
	ID       uint       `gorm:"primarykey"`
	Name     string     `gorm:"unique_index;"`
	Profiles []*Profile `gorm:"many2many:profile_languages;"`
}
