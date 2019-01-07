package invitation

import (
	"couchsport/api/stores"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type InvitationHandler struct {
	Store *stores.InvitationStore
}

func (app InvitationHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()

	users := app.Store.GetInvitations(keys)
	json, err := json.Marshal(users)

	if err != nil {
		log.Error(err)
	}

	fmt.Fprintf(w, string(json))

}
