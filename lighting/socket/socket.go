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

func Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	mu := &sync.Mutex{}

	unsubscribeFn := store.Subscribe(notifyValueChanged(mu, c))
	defer unsubscribeFn()

	for {
		performWebsocketCycle(mu, c)
	}
}
func performWebsocketCycle(mu *sync.Mutex, c *websocket.Conn) {
	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}

	if mt == websocket.TextMessage {
		var details socketPayload
		json.Unmarshal(message, &details)

		switch details.Type {
		case ping:
			log.Println("ping")
		case channelState:
			err = processChannelState(mu, c)
			if err != nil {
				log.Println("channelState:", err)
			}
		case trackDetails:
			log.Println("trackDetails")
		case updateChannel:
			err = processUpdateChannel(mu, c, message)
			if err != nil {
				log.Println("updateChannel:", err)
			}
		}
	}
}



