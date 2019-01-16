// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stores

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type query struct {
	ID              uint
	Action, Message string
}

//hub maintains the set of active clients and broadcasts messages to the
// clients.
type hub struct {
	// Registered clients.
	clients map[uint]*client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Inbound messages from the clients.
	dispatch chan query

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client
}

func newHub() *hub {
	return &hub{
		broadcast:  make(chan []byte),
		dispatch:   make(chan query),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[uint]*client),
	}
}

func (me *hub) run() {
	for {
		select {
		case q := <-me.dispatch:
			me.handleQueries(q)
		case client := <-me.register:
			me.clients[client.ID] = client
		case client := <-me.unregister:
			if _, ok := me.clients[client.ID]; ok {
				delete(me.clients, client.ID)
				close(client.send)
			}
		case message := <-me.broadcast:
			for _, client := range me.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(me.clients, client.ID)
				}
			}
		}
	}
}

func (me *hub) Register(profileID uint, conn *websocket.Conn) {
	client := &client{ID: profileID, hub: me, conn: conn, send: make(chan []byte, 256)}
	me.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (me *hub) handleQueries(q query) {
	log.Printf("ws hub: received query : action = %s, message = %s", q.Action, q.Message)

	// me.emit(query{
	// 	ID:      q.ID,
	// 	Action:  "message.new",
	// 	Message: q.Message,
	// })
}

func (me *hub) Emit(profileID uint, action, message string) {
	log.Printf("ws hub: sending to %d action %s", profileID, action)
	q := query{
		ID:      profileID,
		Action:  action,
		Message: message,
	}

	if err := me.emit(q); err != nil {
		log.Printf("ws hub error: %s", err)
	}
}

func (me *hub) emit(q query) error {
	c := me.clients[q.ID]
	if c == nil {
		return fmt.Errorf("%s", "client not connected")
	}

	jsonBody, err := json.Marshal(q)
	if err != nil {
		return err
	}

	c.send <- jsonBody

	return nil
}
