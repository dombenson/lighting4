// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"lighting/channelUpdater"
	"lighting/lights"
	"lighting/store"
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
	SeqNo    int              `json:"seqNo"`
	FadeTime int              `json:"fadeTime"`
}

func (this *socketConnection) processUpdateChannel(message []byte) error {
	var details updateChannelPayload

	err := json.Unmarshal(message, &details)
	if err != nil {
		return err
	}

	currentValue, currentSeqNo := store.GetValueAndSeqNo(details.Data.Channel.Id)

	if currentValue != details.Data.Channel.Value && currentSeqNo <= details.Data.Channel.SeqNo {
		log.Infof("(%d) 'updateChannel' %d -> %d (%d)", this.id, details.Data.Channel.Id, details.Data.Channel.Value, details.Data.Channel.SeqNo)
		channelUpdater.GetChannelUpdater(details.Data.Channel.Id).UpdateValue(details.Data.Channel.Value)
	} else {
		log.Debugf("(%d) 'updateChannel' [ignored] %d -> %d (%d)", this.id, details.Data.Channel.Id, details.Data.Channel.Value, details.Data.Channel.SeqNo)
	}

	return nil
}
