// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package store

import (
	"github.com/op/go-logging"
	"lighting/amqp/lightingControl"
	"lighting/lights"
	"sync"
)

var log = logging.MustGetLogger("store")

var mu *sync.RWMutex
var values map[lights.ChannelNo]lights.Value
var valueSequences map[lights.ChannelNo]int
var lastSeenSeqs map[lights.ChannelNo]int

func init() {
	values = make(map[lights.ChannelNo]lights.Value)
	valueSequences = make(map[lights.ChannelNo]int)
	lastSeenSeqs = make(map[lights.ChannelNo]int)
	mu = &sync.RWMutex{}
}

func Sync() error {
	err := lightingControl.RequestValues()
	if err != nil {
		return err
	}

	return nil
}

func UpdateValue(channel lights.ChannelNo, value lights.Value) error {
	hasChanged := SetValue(channel, value)

	if !hasChanged {
		return nil
	}

	err := lightingControl.SetValue(channel, value)
	if err != nil {
		return err
	}

	return nil
}

func SetValue(channel lights.ChannelNo, value lights.Value) bool {
	hasChanged, seqNo := doSetValue(channel, value)

	if hasChanged {
		notify(channel, value, seqNo)
	}

	return hasChanged
}

func doSetValue(channel lights.ChannelNo, value lights.Value) (bool, int) {
	mu.Lock()
	defer mu.Unlock()

	originalValue := values[channel]
	values[channel] = value

	hasChanged := originalValue != value
	if hasChanged {
		valueSequences[channel]++
		return true, valueSequences[channel]
	}

	return false, 0
}

func Reset() {
	mu.Lock()
	defer mu.Unlock()

	lastSeenSeqs = make(map[lights.ChannelNo]int)

	for channel := lights.ChannelNo(1); channel <= GetLastCommissionedChannel(); channel++ {
		err := lightingControl.SetValue(channel, values[channel])
		if err != nil {
			log.Error("Unable to transmit existing value", err)
		}
	}
}

func GetValue(channel lights.ChannelNo) lights.Value {
	mu.RLock()
	defer mu.RUnlock()

	return values[channel]
}

func GetLastCommissionedChannel() lights.ChannelNo {
	return 512
}

func GetLastSeenSeqNo(channelNo lights.ChannelNo) int {
	mu.RLock()
	defer mu.RUnlock()

	return lastSeenSeqs[channelNo]
}

func SetLastSeenSeqNo(channelNo lights.ChannelNo, seqNo int) {
	mu.Lock()
	defer mu.Unlock()

	lastSeenSeqs[channelNo] = seqNo
}