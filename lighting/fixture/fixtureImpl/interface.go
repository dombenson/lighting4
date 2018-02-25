// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import (
	"lighting/fixtureType"
	"lighting/lights"
	"time"
)

type FixtureImpl interface {
	GetName() string
	GetDescription() string
	GetType() *fixtureType.FixtureType
	GetFirstChannel() lights.ChannelNo
}

type RGBFixtureImpl interface {
	FixtureImpl

	GetRedValue() lights.Value
	GetGreenValue() lights.Value
	GetBlueValue() lights.Value
	SetRedValue(value lights.Value, fade time.Duration)
	SetGreenValue(value lights.Value, fade time.Duration)
	SetBlueValue(value lights.Value, fade time.Duration)
	SetColor(red, green, blue lights.Value, fade time.Duration)
}
