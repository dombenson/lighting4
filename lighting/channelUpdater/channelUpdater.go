// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package channelUpdater

import (
	"lighting/lights"
	"lighting/store"
	"lighting/util/debouncer"
	"log"
	"time"
)

var channelUpdaters map[lights.ChannelNo]*ChannelUpdater

func init() {
	channelUpdaters = make(map[lights.ChannelNo]*ChannelUpdater)
}

type ChannelUpdater struct {
	channelNo        lights.ChannelNo
	channelDebouncer *debouncer.Debouncer
}

func GetChannelUpdater(channelNo lights.ChannelNo) *ChannelUpdater {
	updater, ok := channelUpdaters[channelNo]

	if updater == nil || !ok {
		updater = newChannelUpdater(channelNo)
		channelUpdaters[channelNo] = updater
	}

	return updater
}

func newChannelUpdater(channelNo lights.ChannelNo) *ChannelUpdater {
	updater := ChannelUpdater{
		channelNo: channelNo,
	}

	channelDebouncer := debouncer.New(debouncer.Opts{
		Duration: 20 * time.Millisecond,
		Callback: func(data interface{}) {
			value, ok := data.(lights.Value)
			if !ok {
				log.Printf("[channelUpdate] (%d) type assertion failed\n", channelNo)
				return
			}

			err := store.UpdateValue(channelNo, value)
			if err != nil {
				log.Printf("[channelUpdate] (%d) error updating value (%s)", channelNo, err)
				return
			}

			log.Printf("[channelUpdate] (%d) update value to %d\n", channelNo, value)
		},
	})

	updater.channelDebouncer = channelDebouncer

	return &updater
}

func (this *ChannelUpdater) UpdateValue(value lights.Value) {
	this.channelDebouncer.Set(value)
}
