// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package store

import (
	"lighting/lights"
	"sync"
)

var observers []chan ValuesChange
var mu sync.Mutex

type ValuesChange struct {
	Channel lights.ChannelNo
	Value   lights.Value
}

func Attach(c chan ValuesChange) {
	mu.Lock()
	defer mu.Unlock()
	observers = append(observers, c)
}

func Detach(c chan ValuesChange) {
	mu.Lock()
	defer mu.Unlock()
	for i, v := range observers {
		if v == c {
			observers = append(observers[:i], observers[i+1:]...)
			return
		}
	}
}

func notify(channel lights.ChannelNo, value lights.Value) {
	for _, v := range observers {
		v <- ValuesChange {
			Channel: channel,
			Value: value,
		}
	}
}
