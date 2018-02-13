// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package store

import (
	"lighting/amqp/lightingControl"
	"lighting/lights"
)

var values map[lights.ChannelNo]lights.Value

func init() {
	values = make(map[lights.ChannelNo]lights.Value)
}

func UpdateValue(channel lights.ChannelNo, value lights.Value) error {
	SetValue(channel, value)

	err := lightingControl.SetValue(channel, value)
	if err != nil {
		return err
	}

	return nil
}

func SetValue(channel lights.ChannelNo, value lights.Value) {
	originalValue := values[channel]
	values[channel] = value
	if originalValue != value {
		notify(channel, value)
	}
}

func GetValue(channel lights.ChannelNo) lights.Value {
	return values[channel]
}