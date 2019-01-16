package stores

import (
	"github.com/goland-amaurybrisou/couchsport/api/types"
	"github.com/goland-amaurybrisou/couchsport/config"
	"github.com/jinzhu/gorm"
)

//StoreFactory holds references to every Store in the application
type StoreFactory struct {
	Db                *gorm.DB
	wsStore           *hub
	activityStore     *activityStore
	languageStore     *languageStore
	imageStore        *imageStore
	userStore         *userStore
	sessionStore      *sessionStore
	fileStore         *fileStore
	profileStore      *profileStore
	pageStore         *pageStore
	conversationStore *conversationStore
}

//NewStoreFactory is the first store layer. ask him what store you want
func NewStoreFactory(Db *gorm.DB, c config.Config) *StoreFactory {

	hub := newHub()

	fileStore := fileStore{
		FileSystem:    types.OsFS{},
		PublicPath:    c.PublicPath,
		ImageBasePath: c.ImageBasePath,
		FilePrefix:    c.FilePrefix,
	}

	profileStore := profileStore{Db: Db, FileStore: fileStore}

	return &StoreFactory{
		wsStore:           hub,
		activityStore:     &activityStore{Db: Db},
		languageStore:     &languageStore{Db: Db},
		imageStore:        &imageStore{Db: Db},
		userStore:         &userStore{Db: Db},
		sessionStore:      &sessionStore{Db: Db},
		fileStore:         &fileStore,
		profileStore:      &profileStore,
		pageStore:         &pageStore{Db: Db, FileStore: fileStore, ProfileStore: profileStore},
		conversationStore: &conversationStore{Db: Db},
	}
}

//Init initialize Databse tables
func (me StoreFactory) Init(populate bool) {
	go me.wsStore.run()

	if !populate {
		return
	}
	me.profileStore.Migrate()

	me.userStore.Migrate() //profile needs profile

	me.sessionStore.Migrate() //session needs user

	me.languageStore.Migrate()     //language need profile
	me.pageStore.Migrate()         //page needs profile
	me.conversationStore.Migrate() //conversation needs profile

	me.activityStore.Migrate() //activity needs page & profile
	me.imageStore.Migrate()    //image needs page

}

//WsStore returns the app wesocket hub
func (me StoreFactory) WsStore() *hub {
	return me.wsStore
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

//ConversationStore returns the app userStore
func (me StoreFactory) ConversationStore() *conversationStore {
	return me.conversationStore
}
