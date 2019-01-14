package handlers

import (
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/stores"
	"net/http"
)

type conversationHandler struct {
	Store *stores.StoreFactory
}

func (me conversationHandler) HandleMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not yet Implemented")
}

func (me conversationHandler) ProfileConversations(userID uint, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not yet Implemented")
}

func (me conversationHandler) Delete(userID uint, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Not yet Implemented")
}
