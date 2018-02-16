// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package store

import (
	"github.com/satori/go.uuid"
	"lighting/lights"
)

var subscribers map[uuid.UUID]ValueChangeCallback

type ValuesChange struct {
	Channel lights.ChannelNo
	Value   lights.Value
	SeqNo   int
}

type ValueChangeCallback func(change ValuesChange)

func init() {
	subscribers = make(map[uuid.UUID]ValueChangeCallback)
}

func Subscribe(callback ValueChangeCallback) func() {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	subscribers[id] = callback

	return func() {
		delete(subscribers, id)
	}
}

func notify(channel lights.ChannelNo, value lights.Value, seqNo int) {
	for _, callback := range subscribers {
		callback(ValuesChange {
			Channel: channel,
			Value:   value,
			SeqNo:   seqNo,
		})
	}
}
