package handlers

import (
	"fmt"
	"github.com/goland-amaurybrisou/couchsport/api/stores"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

//wsHandler ...
type wsHandler struct {
	WsUpgrader *websocket.Upgrader
	Stores     *stores.StoreFactory
}

//Emit message
// func (me *wsHandler) Emit(profileID uint, action, message string) {
// 	query := query{
// 		ID:      profileID,
// 		Action:  action,
// 		Message: message,
// 	}
// 	err := me.Hub.emit(query)
// 	if err != nil {
// 		log.Printf("ws hub error : %s", err)
// 	}
// }

//EntryPoint Ws handler
func (me *wsHandler) EntryPoint(w http.ResponseWriter, r *http.Request) {
	stringID := r.URL.Query().Get("id")
	if stringID == "" {
		log.Println("ws: invalid id")
		http.Error(w, fmt.Errorf("ws: invalid id %s", stringID).Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(stringID)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Errorf("%s", err).Error(), http.StatusBadRequest)
		return
	}

	conn, err := me.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	me.Stores.WsStore().Register(uint(id), conn)
}

func (me *wsHandler) echo(conn *websocket.Conn, mt int, message []byte) {
	err := conn.WriteMessage(mt, message)
	if err != nil {
		log.Println("write:", err)
	}
}
