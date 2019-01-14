package stores

import (
	"github.com/goland-amaurybrisou/couchsport/api/models"
	"github.com/jinzhu/gorm"
)

type languageStore struct {
	Db *gorm.DB
}

//Migrate creates the db table
func (me languageStore) Migrate() {
	me.Db.AutoMigrate(&models.Language{})
	me.Db.Table("profile_languages").AddForeignKey("language_id", "languages(id)", "NO ACTION", "NO ACTION")
	me.Db.Table("profile_languages").AddForeignKey("profile_id", "profiles(id)", "CASCADE", "NO ACTION")

	me.Db.Table("profile_languages").AddUniqueIndex("language_id_profile_id_unique", "profile_id, language_id")

	languages := []string{"Abkhazian", "Afar", "Afrikaans", "Akan", "Albanian", "Amharic", "Arabic", "Aragonese", "Armenian", "Assamese", "Avar", "Aymara", "Azerbaijani", "Bambara", "Bashkir", "Basque", "Belarusian", "Bengali", "Bihari", "Bislama", "Bosnian", "Breton", "Bulgarian", "Burmese", "Cambodian", "Catalan", "Chamorro", "Chechen", "Chichewa", "Chinese", "Chuvash", "Cornish", "Corsican", "Cree", "Croatian", "Czech", "Danish", "Divehi", "Dutch", "Dzongkha", "English", "Esperanto", "Estonian", "Ewe", "Faroese", "Fijian", "Finnish", "French", "Galician", "Ganda", "Georgian", "German", "Greek", "Greenlandic", "Guarani", "Gujarati", "Haitian", "Hausa", "Hebrew", "Herero", "Hindi", "Hiri Motu", "Hungarian", "Icelandic", "Ido", "Igbo", "Indonesian", "Interlingua", "Interlingue", "Inuktitut", "Inupiak", "Irish", "Italian", "Japanese", "Javanese", "Kannada", "Kanuri", "Kashmiri", "Kazakh", "Kikuyu", "Kirghiz", "Kirundi", "Komi", "Kongo", "Korean", "Kuanyama", "Kurdish", "Laotian", "Latin", "Latvian", "Limburgian", "Lingala", "Lithuanian", "Luxembourgish", "Macedonian", "Malagasy", "Malay", "Malayalam", "Maltese", "Manx", "Maori", "Marathi", "Marshallese", "Moldovan", "Mongolian", "Nauruan", "Navajo", "Ndonga", "Nepali", "North Ndebele", "Northern Sami", "Norwegian", "Norwegian Nynorsk", "Occitan", "Ojibwa", "Oriya", "Oromo", "Ossetian / Ossetic", "Pali", "Panjabi / Punjabi", "Pashto", "Persian", "Peul", "Polish", "Portuguese", "Quechua", "Raeto Romance", "Romanian", "Russian", "Rwandi", "Samoan", "Sango", "Sanskrit", "Sardinian", "Scottish Gaelic", "Serbian", "Serbo-Croatian", "Shona", "Sichuan Yi", "Sindhi", "Sinhalese", "Slovak", "Slovenian", "Somalia", "South Ndebele", "Southern Sotho", "Spanish", "Sundanese", "Swahili", "Swati", "Swedish", "Tagalog / Filipino", "Tahitian", "Tajik", "Tamil", "Tatar", "Telugu", "Thai", "Tibetan", "Tigrinya", "Tonga", "Tsonga", "Tswana", "Turkish", "Turkmen", "Twi", "Ukrainian", "Urdu", "Uyghur", "Uzbek", "Venda", "Vietnamese", "Volapük", "Walloon", "Welsh", "West Frisian", "Wolof", "Xhosa", "Yiddish", "Yoruba", "Zhuang", "Zulu"}

	for _, l := range languages {
		me.Db.FirstOrCreate(&models.Language{Name: l}, models.Language{Name: l})
	}
}

//All returns all the languages in db
func (me languageStore) All() ([]models.Language, error) {
	var languages []models.Language
	if err := me.Db.Find(&languages).Error; err != nil {
		return []models.Language{}, err
	}
	return languages, nil
}
