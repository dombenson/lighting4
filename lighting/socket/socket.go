// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"lighting/store"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{} // use default options
var lastId int

type socketPayload struct {
	Type string `json:"type"`
}

const (
	channelState  = "channelState"
	trackDetails  = "trackDetails"
	updateChannel = "updateChannel"
	ping          = "ping"
	notifyChange  = "uC"
)

type socketConnection struct{
	id            int
	mu            *sync.Mutex
	c             *websocket.Conn
	unsubscribeFn func()
}

func newSocketConnection(w http.ResponseWriter, r *http.Request) (*socketConnection, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	nextId := lastId + 1
	lastId = nextId

	connection := &socketConnection{
		id: nextId,
		mu: &sync.Mutex{},
		c:  c,
	}

	connection.unsubscribeFn = store.Subscribe(connection.notifyValueChanged())

	return connection, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	connection, err := newSocketConnection(w, r)
	if err != nil {
		log.Printf("[socket] could not open web socket connection (%s)\n", err)
		return
	}
	defer connection.close()

	log.Printf("[socket] (%d) connected from '%s'\n", connection.id, r.RemoteAddr)

	connection.start()
}

func (this *socketConnection) start() {
	for {
		hasClosed := this.performWebsocketCycle()
		if hasClosed {
			break
		}
	}
}

func (this *socketConnection) close() error {
	this.unsubscribeFn()

	err := this.c.Close()
	if err != nil {
		return err
	}

	return nil
}

func (this *socketConnection) performWebsocketCycle() bool {
	mt, message, err := this.c.ReadMessage()
	if err != nil {

		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
			log.Printf("[socket] (%d) closed\n", this.id)
			return true
		}

		log.Printf("[socket] (%d) error (%s)\n", this.id, err)
		return true
	}

	if mt == websocket.TextMessage {
		var details socketPayload
		json.Unmarshal(message, &details)

		switch details.Type {
		case ping:
			log.Printf("[socket] (%d) 'ping'\n", this.id)
		case channelState:
			err = this.processChannelState()
			if err != nil {
				log.Printf("[socket] (%d) 'channelState' processing error (%s)\n", this.id, err)
			}
		case trackDetails:
			log.Println("trackDetails")
		case updateChannel:
			err = this.processUpdateChannel(message)
			if err != nil {
				log.Printf("[socket] (%d) 'updateChannel' processing error (%s)\n", this.id, err)
			}
		}
	}

	return false
}



