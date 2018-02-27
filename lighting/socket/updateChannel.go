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
	Universe int              `json:"universe"`
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

	address := lights.NewAddress(details.Data.Channel.Universe, details.Data.Channel.Id)

	currentValue, currentSeqNo := store.GetValueAndSeqNo(address)

	if currentValue != details.Data.Channel.Value && currentSeqNo <= details.Data.Channel.SeqNo {
		log.Infof("(%d) 'updateChannel' %d:%d -> %d (%d)", this.id, address.Universe, address.ChannelNo, details.Data.Channel.Value, details.Data.Channel.SeqNo)
		channelUpdater.GetChannelUpdater(address).UpdateValue(details.Data.Channel.Value)
	} else {
		log.Debugf("(%d) 'updateChannel' [ignored] %d:%d -> %d (%d)", this.id, address.Universe, address.ChannelNo, details.Data.Channel.Value, details.Data.Channel.SeqNo)
	}

	return nil
}
