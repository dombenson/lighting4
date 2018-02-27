// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
	"lighting/store"
	"net/http"
	"sync"
)

var log = logging.MustGetLogger("socket")

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
		log.Errorf("could not open web socket connection (%s)", err)
		return
	}
	defer connection.close()

	log.Infof("(%d) connected from '%s'", connection.id, r.RemoteAddr)

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
			log.Infof("(%d) closed", this.id)
			return true
		}

		log.Errorf("(%d) error (%s)", this.id, err)
		return true
	}

	if mt == websocket.TextMessage {
		var details socketPayload
		json.Unmarshal(message, &details)

		switch details.Type {
		case ping:
			log.Debugf("(%d) 'ping'", this.id)
		case channelState:
			err = this.processChannelState()
			if err != nil {
				log.Errorf("(%d) 'channelState' processing error (%s)", this.id, err)
			}
		case trackDetails:
			log.Infof("(%d) 'trackDetails' not currently handled", this.id)
		case updateChannel:
			err = this.processUpdateChannel(message)
			if err != nil {
				log.Errorf("(%d) 'updateChannel' processing error (%s)", this.id, err)
			}
		}
	}

	return false
}



