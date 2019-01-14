package handlers

import (
	"github.com/goland-amaurybrisou/couchsport/api/stores"
)

//HandlerFactory hols all the handler of the application
type HandlerFactory struct {
	activityHandler     activityHandler
	imageHandler        imageHandler
	languageHandler     languageHandler
	pageHandler         pageHandler
	profileHandler      profileHandler
	userHandler         userHandler
	conversationHandler conversationHandler
}

//NewHandlerFactory generates the handlerFactory holding every handler in the application
func NewHandlerFactory(storeFactory *stores.StoreFactory) *HandlerFactory {
	return &HandlerFactory{
		activityHandler:     activityHandler{Stores: storeFactory},
		imageHandler:        imageHandler{Store: storeFactory},
		languageHandler:     languageHandler{Store: storeFactory},
		pageHandler:         pageHandler{Store: storeFactory},
		profileHandler:      profileHandler{Store: storeFactory},
		userHandler:         userHandler{Store: storeFactory},
		conversationHandler: conversationHandler{Store: storeFactory},
	}
}

//ActivityHandler returns the applicatioin ActivityHandler
func (me HandlerFactory) ActivityHandler() *activityHandler {
	return &me.activityHandler
}

//ImageHandler returns the applicatioin ImageHandler
func (me HandlerFactory) ImageHandler() *imageHandler {
	return &me.imageHandler
}

//LanguageHandler returns the applicatioin LanguageHandler
func (me HandlerFactory) LanguageHandler() *languageHandler {
	return &me.languageHandler
}

//PageHandler returns the applicatioin PageHandler
func (me HandlerFactory) PageHandler() *pageHandler {
	return &me.pageHandler
}

//ProfileHandler returns the applicatioin ProfileHandler
func (me HandlerFactory) ProfileHandler() *profileHandler {
	return &me.profileHandler
}

//UserHandler returns the applicatioin UserHandler
func (me HandlerFactory) UserHandler() *userHandler {
	return &me.userHandler
}

//ConversationHandler returns the applicatioin ConversationHandler
func (me HandlerFactory) ConversationHandler() *conversationHandler {
	return &me.conversationHandler
}
