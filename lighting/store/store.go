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
var values map[lights.Address]lights.Value
var valueSequences map[lights.Address]int
var lastSeenHardwareSeqNos map[lights.Address]int

func init() {
	values = make(map[lights.Address]lights.Value)
	valueSequences = make(map[lights.Address]int)
	lastSeenHardwareSeqNos = make(map[lights.Address]int)
	mu = &sync.RWMutex{}
}

func Sync() error {
	for _, universe := range GetUniverses() {
		err := lightingControl.RequestValues(universe)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateValue(channel lights.Address, value lights.Value) error {
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

func SetValue(channel lights.Address, value lights.Value) bool {
	hasChanged, seqNo := doSetValue(channel, value)

	if hasChanged {
		notify(channel, value, seqNo)
	}

	return hasChanged
}

func doSetValue(channel lights.Address, value lights.Value) (bool, int) {
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

func Reset(universe int) {
	mu.Lock()
	defer mu.Unlock()

	lastSeenHardwareSeqNos = make(map[lights.Address]int)

	for channel := lights.ChannelNo(1); channel <= GetLastCommissionedChannel(universe); channel++ {
		address := lights.NewAddress(universe, channel)
		err := lightingControl.SetValue(address, values[address])
		if err != nil {
			log.Error("Unable to transmit existing value", err)
		}
	}
}

func GetValue(channel lights.Address) lights.Value {
	mu.RLock()
	defer mu.RUnlock()

	return values[channel]
}

func GetValueAndSeqNo(channel lights.Address) (lights.Value, int) {
	mu.RLock()
	defer mu.RUnlock()

	return values[channel], valueSequences[channel]
}

func GetSeqNo(channel lights.Address) int {
	mu.RLock()
	defer mu.RUnlock()

	return valueSequences[channel]
}

func GetUniverses() []int {
	return []int{1, 2}
}

func GetLastCommissionedChannel(universe int) lights.ChannelNo {
	return 512
}

func GetLastSeenHardwareSeqNo(channelNo lights.Address) int {
	mu.RLock()
	defer mu.RUnlock()

	return lastSeenHardwareSeqNos[channelNo]
}

func SetLastSeenHardwareSeqNo(channelNo lights.Address, seqNo int) {
	mu.Lock()
	defer mu.Unlock()

	lastSeenHardwareSeqNos[channelNo] = seqNo
}