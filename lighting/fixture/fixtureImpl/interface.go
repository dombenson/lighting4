// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import (
	"github.com/brutella/hc/accessory"
	"lighting/fixtureType"
	"lighting/lights"
	"time"
)

type FixtureImpl interface {
	GetName() string
	GetDescription() string
	GetType() *fixtureType.FixtureType
	GetFirstChannel() lights.Address
}

type RGBFixtureImpl interface {
	FixtureImpl

	GetHomeKitAccessory() *accessory.Accessory

	IsFaderAvailable() bool
	IsWhiteAvailable() bool
	IsAmberAvailable() bool
	IsUvAvailable() bool
	GetFaderValue() lights.Value
	GetRedValue() lights.Value
	GetGreenValue() lights.Value
	GetBlueValue() lights.Value
	GetWhiteValue() lights.Value
	GetAmberValue() lights.Value
	GetUvValue() lights.Value
	SetFaderValue(value lights.Value, fade time.Duration)
	SetRedValue(value lights.Value, fade time.Duration)
	SetGreenValue(value lights.Value, fade time.Duration)
	SetBlueValue(value lights.Value, fade time.Duration)
	SetWhiteValue(value lights.Value, fade time.Duration)
	SetAmberValue(value lights.Value, fade time.Duration)
	SetUvValue(value lights.Value, fade time.Duration)
	SetColor(red, green, blue lights.Value, fade time.Duration)
}
