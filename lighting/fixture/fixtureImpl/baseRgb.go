// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import (
	"lighting/channelUpdater"
	"lighting/lights"
	"lighting/store"
	"time"
)

type baseRGBFixture struct {
	redChannel   lights.ChannelNo
	greenChannel lights.ChannelNo
	blueChannel  lights.ChannelNo
}

func (this *baseRGBFixture) GetRedValue() lights.Value {
	return store.GetValue(this.redChannel)
}

func (this *baseRGBFixture) GetGreenValue() lights.Value {
	return store.GetValue(this.greenChannel)
}

func (this *baseRGBFixture) GetBlueValue() lights.Value {
	return store.GetValue(this.blueChannel)
}

func (this *baseRGBFixture) SetRedValue(value lights.Value, fade time.Duration) {
	channelUpdater.GetChannelUpdater(this.redChannel).UpdateValueWithFade(this.GetRedValue(), value, fade)
}

func (this *baseRGBFixture) SetGreenValue(value lights.Value, fade time.Duration) {
	channelUpdater.GetChannelUpdater(this.greenChannel).UpdateValueWithFade(this.GetGreenValue(), value, fade)
}

func (this *baseRGBFixture) SetBlueValue(value lights.Value, fade time.Duration) {
	channelUpdater.GetChannelUpdater(this.blueChannel).UpdateValueWithFade(this.GetBlueValue(), value, fade)
}

func (this *baseRGBFixture) SetColor(red, green, blue lights.Value, fade time.Duration) {
	this.SetRedValue(red, fade)
	this.SetGreenValue(green, fade)
	this.SetBlueValue(blue, fade)
}