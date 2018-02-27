// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package channelUpdater

import (
	"github.com/op/go-logging"
	"lighting/lights"
	"lighting/store"
	"lighting/util/debouncer"
	"reflect"
	"time"
)

var log = logging.MustGetLogger("channelUpdater")

var channelUpdaters = make(map[lights.Address]*ChannelUpdater)

type ChannelUpdater struct {
	address          lights.Address
	channelDebouncer *debouncer.Debouncer
}

func GetChannelUpdater(address lights.Address) *ChannelUpdater {
	updater, ok := channelUpdaters[address]

	if updater == nil || !ok {
		updater = newChannelUpdater(address)
		channelUpdaters[address] = updater
	}

	return updater
}

func newChannelUpdater(address lights.Address) *ChannelUpdater {
	updater := ChannelUpdater{
		address: lights.Address{},
	}

	channelDebouncer := debouncer.New(debouncer.Opts{
		Duration: 20 * time.Millisecond,
		Callback: func(data interface{}) {
			value, ok := data.(lights.Value)
			if !ok {
				log.Errorf("(%d:%d) type assertion failed (%s)", address.Universe, address.ChannelNo, reflect.TypeOf(data))
				return
			}

			err := store.UpdateValue(address, value)
			if err != nil {
				log.Errorf("(%d:%d) error updating value (%s)", address.Universe, address.ChannelNo, err)
				return
			}
		},
	})

	updater.channelDebouncer = channelDebouncer

	return &updater
}

func (this *ChannelUpdater) UpdateValueWithFade(startValue, endValue lights.Value, duration time.Duration) {
	if duration == 0 {
		this.UpdateValue(endValue)
		return
	}

	this.UpdateValue(startValue)

	if startValue == endValue {
		return
	}

	stepDuration := calculateStepDuration(startValue, endValue, duration)

	ticker := time.NewTicker(stepDuration)

	currentValue := startValue

	go func() {
		for range ticker.C {
			if startValue < endValue {
				currentValue += 1
			} else {
				currentValue -= 1
			}

			this.UpdateValue(currentValue)

			if currentValue == endValue {
				ticker.Stop()
			}
		}
	}()
}

func calculateStepDuration(startValue, endValue lights.Value, duration time.Duration) time.Duration {
	if startValue < endValue {
		return time.Duration(int64(duration) / int64(endValue - startValue))
	} else {
		return time.Duration(int64(duration) / int64(startValue - endValue))
	}
}

func (this *ChannelUpdater) UpdateValue(value lights.Value) {
	this.channelDebouncer.Set(value)
}
