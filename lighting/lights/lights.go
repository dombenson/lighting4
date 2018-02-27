// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package lights

type Address struct {
	Universe  int       `json:"universe"`
	ChannelNo ChannelNo `json:"channel"`
}

func NewAddress(universe int, channelNo ChannelNo) Address {
	if universe == 0 {
		universe = 1
	}

	return Address{
		Universe:  universe,
		ChannelNo: channelNo,
	}
}

type ChannelNo uint16
type Value uint8
