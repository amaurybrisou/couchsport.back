package stores

import (
	"couchsport/api/types"
	"couchsport/config"
	"github.com/jinzhu/gorm"
)

//StoreFactory holds references to every Store in the application
type StoreFactory struct {
	Db                 *gorm.DB
	activityStore      *activityStore
	languageStore      *languageStore
	imageStore         *imageStore
	userStore          *userStore
	sessionStore       *sessionStore
	fileStore          *fileStore
	profileStore       *profileStore
	pageStore          *pageStore
	authorizationStore *authorizationStore
}

//NewStoreFactory is the first store layer. ask him what store you want
func NewStoreFactory(Db *gorm.DB, c config.Config) *StoreFactory {
	fileStore := fileStore{
		FileSystem:    types.OsFS{},
		PublicPath:    c.PublicPath,
		ImageBasePath: c.ImageBasePath,
		FilePrefix:    c.FilePrefix,
	}

	profileStore := profileStore{Db: Db, FileStore: fileStore}

	return &StoreFactory{
		activityStore:      &activityStore{Db: Db},
		languageStore:      &languageStore{Db: Db},
		imageStore:         &imageStore{Db: Db},
		userStore:          &userStore{Db: Db},
		sessionStore:       &sessionStore{Db: Db},
		fileStore:          &fileStore,
		profileStore:       &profileStore,
		pageStore:          &pageStore{Db: Db, FileStore: fileStore, ProfileStore: profileStore},
		authorizationStore: &authorizationStore{Db: Db},
	}
}

//Init initialize Databse tables
func (me StoreFactory) Init() {
	me.activityStore.Migrate()
	me.languageStore.Migrate()
	me.imageStore.Migrate()
	me.userStore.Migrate()
	me.sessionStore.Migrate()
	me.profileStore.Migrate()
	me.pageStore.Migrate()
}

//AuthorizationStore returns the app authorizationStore
func (me StoreFactory) AuthorizationStore() *authorizationStore {
	return me.authorizationStore
}

//PageStore returns the app pageStore
func (me StoreFactory) PageStore() *pageStore {
	return me.pageStore
}

//FileStore returns the app fileStore
func (me StoreFactory) FileStore() *fileStore {
	return me.fileStore
}

//ImageStore returns the app imageStore
func (me StoreFactory) ImageStore() *imageStore {
	return me.imageStore
}

//ProfileStore returns the app profileStore
func (me StoreFactory) ProfileStore() *profileStore {
	return me.profileStore
}

//SessionStore returns the app sessionStore
func (me StoreFactory) SessionStore() *sessionStore {
	return me.sessionStore
}

//LanguageStore returns the app languageStore
func (me StoreFactory) LanguageStore() *languageStore {
	return me.languageStore
}

//ActivityStore returns the app activityStore
func (me StoreFactory) ActivityStore() *activityStore {
	return me.activityStore
}

//UserStore returns the app userStore
func (me StoreFactory) UserStore() *userStore {
	return me.userStore
}
