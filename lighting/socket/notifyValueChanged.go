// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"lighting/lights"
	"lighting/store"
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
	Universe  int              `json:"u"`
	Value     lights.Value     `json:"l"`
	SeqNo     int              `json:"s"`
}

func (this *socketConnection) notifyValueChanged() store.ValueChangeCallback {
	return func(change store.ValuesChange) {
		this.mu.Lock()
		defer this.mu.Unlock()

		details := notifyChangePayload {
			socketPayload: socketPayload {notifyChange},
			Data: notifyChangeData{
				Channel: notifyChangeDetails{
					ChannelNo: change.Channel.ChannelNo,
					Universe:  change.Channel.Universe,
					Value:     change.Value,
					SeqNo:     change.SeqNo,
				},
			},
		}

		message, err := json.Marshal(details)
		if err != nil {
			log.Errorf("(%d:%d) Error in encoding notify (%s)", this.id, err)
			return
		}

		err = this.c.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Errorf("(%d:%d) Error in notify (%s)", this.id, err)
			return
		}

		return
	}
}
