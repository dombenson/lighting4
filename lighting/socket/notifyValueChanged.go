// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"lighting/lights"
	"lighting/store"
	"log"
	"sync"
)

type notifyChangePayload struct {
	socketPayload
	Data notifyChangeData `json:"data"`
}

type notifyChangeData struct {
	Channel notifyChangeDetails `json:"c"`
}

type notifyChangeDetails struct {
	ChannelNo lights.ChannelNo `json:"i"`
	Value     lights.Value     `json:"l"`
	SeqNo     int              `json:"s"`
}

func notifyValueChanged(mu *sync.Mutex, c *websocket.Conn) store.ValueChangeCallback {
	return func(change store.ValuesChange) {
		mu.Lock()
		defer mu.Unlock()

		details := notifyChangePayload {
			socketPayload: socketPayload {notifyChange},
			Data: notifyChangeData{
				Channel: notifyChangeDetails{
					ChannelNo: change.Channel,
					Value:     change.Value,
					SeqNo:     change.SeqNo,
				},
			},
		}

		message, err := json.Marshal(details)
		if err != nil {
			log.Println("Error in notify", err)
			return
		}

		err = c.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error in notify", err)
			return
		}

		return
	}
}
