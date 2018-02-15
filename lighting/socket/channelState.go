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

type channelStatePayload struct {
	socketPayload
	Data channelStateData `json:"data"`
}

type channelStateData struct {
	Channels []channelStateValue `json:"channels"`
}

type channelStateValue struct {
	Id           lights.ChannelNo `json:"id"`
	CurrentLevel lights.Value     `json:"currentLevel"`
}

func processChannelState(mu *sync.Mutex, c *websocket.Conn) error {
	log.Println("channelState")

	lastCommissionedChannel := store.GetLastCommissionedChannel()

	channelStates := make([]channelStateValue, 0, lastCommissionedChannel)

	for i := lights.ChannelNo(1); i < lastCommissionedChannel; i++ {
		channelStates = append(channelStates, channelStateValue{
			Id:           i,
			CurrentLevel: store.GetValue(i),
		})
	}

	mu.Lock()
	defer mu.Unlock()

	details := channelStatePayload {
		socketPayload: socketPayload {channelState},
		Data: channelStateData {
			Channels: channelStates,
		},
	}

	message, err := json.Marshal(details)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return err
}
