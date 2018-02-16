// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"lighting/channelUpdater"
	"lighting/lights"
	"log"
)

type updateChannelPayload struct {
	socketPayload
	Data updateChannelData `json:"data"`
}

type updateChannelData struct {
	Channel updateChannelValue `json:"channel"`
}

type updateChannelValue struct {
	Id       lights.ChannelNo `json:"id"`
	Value    lights.Value     `json:"level"`
	FadeTime int              `json:"fadeTime"`
}

func (this *socketConnection) processUpdateChannel(message []byte) error {
	var details updateChannelPayload

	err := json.Unmarshal(message, &details)
	if err != nil {
		return err
	}

	log.Printf("[socket] (%d) 'uC' %d -> %d\n", this.id, details.Data.Channel.Id, details.Data.Channel.Value)

	channelUpdater.GetChannelUpdater(details.Data.Channel.Id).UpdateValue(details.Data.Channel.Value)

	return nil
}
