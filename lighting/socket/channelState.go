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
	SeqNo        int              `json:"seqNo"`
}

func (this *socketConnection) processChannelState() error {
	lastCommissionedChannel := store.GetLastCommissionedChannel()

	channelStates := make([]channelStateValue, 0, lastCommissionedChannel)

	for i := lights.ChannelNo(1); i <= lastCommissionedChannel; i++ {
		currentValue, currentSeqNo := store.GetValueAndSeqNo(i)

		channelStates = append(channelStates, channelStateValue{
			Id:           i,
			CurrentLevel: currentValue,
			SeqNo:        currentSeqNo,
		})
	}

	this.mu.Lock()
	defer this.mu.Unlock()

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

	err = this.c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	log.Infof("(%d) 'channelState' sent", this.id)

	return err
}
