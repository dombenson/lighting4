// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import (
	"channellight"
	hkAccessory "github.com/brutella/hc/accessory"
	hkService "github.com/brutella/hc/service"
	"github.com/op/go-logging"
	"lighting/channelUpdater"
	"lighting/lights"
	"lighting/store"
	"time"
)

var log = logging.MustGetLogger("baseRGBFixture")

const channelNotSpecified lights.ChannelNo = 0

type baseRGBFixture struct {
	colorFixtureChannels colorFixtureChannels
	lightbulb            *hkService.Lightbulb
	accessory            *hkAccessory.Accessory
	lightModel           channellight.ChannelLight
}

type colorFixtureChannels struct {
	fader lights.ChannelNo
	red   lights.ChannelNo
	green lights.ChannelNo
	blue  lights.ChannelNo
	white lights.ChannelNo
	amber lights.ChannelNo
	uv    lights.ChannelNo
}

func newBaseRGBFixture(name string, channels colorFixtureChannels) *baseRGBFixture {
	lightbulb := hkService.NewLightbulb()
	accessory := hkAccessory.New(hkAccessory.Info{Name: name}, hkAccessory.TypeLightbulb)

	accessory.AddService(lightbulb.Service)

	var lightModel channellight.ChannelLight

	if channels.red != channelNotSpecified &&
	   channels.green != channelNotSpecified &&
	   channels.blue != channelNotSpecified {
		if channels.white != channelNotSpecified {
			lightModel = &channellight.SevenChannelLight{}
		} else {
			lightModel = &channellight.FourChannelLight{}
		}
	} else {
		panic("attempted to create an RGB fixture without channels specified for red, green and blue")
	}

	baseRGBFixture := &baseRGBFixture{
		colorFixtureChannels: channels,
		lightbulb:            lightbulb,
		accessory:            accessory,
		lightModel:           lightModel,
	}

	lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		baseRGBFixture.syncColorsForLight()
	})

	lightbulb.Brightness.OnValueRemoteUpdate(func(b int) {
		baseRGBFixture.syncColorsForLight()
	})

	lightbulb.Hue.OnValueRemoteUpdate(func(h float64) {
		baseRGBFixture.syncColorsForLight()
	})

	lightbulb.Saturation.OnValueRemoteUpdate(func(s float64) {
		baseRGBFixture.syncColorsForLight()
	})

	store.Subscribe(baseRGBFixture.ValueChange)

	return baseRGBFixture
}

func (this *baseRGBFixture) syncColorsForLight() {
	this.lightModel.SetColor(this.lightbulb)

	switch v := this.lightModel.(type) {
	case *channellight.SevenChannelLight:
		outputColor := v.GetOutputColor()

		this.SetFaderValue(lights.Value(outputColor.Brightness), time.Duration(0))
		this.SetRedValue(lights.Value(outputColor.Red), time.Duration(0))
		this.SetGreenValue(lights.Value(outputColor.Green), time.Duration(0))
		this.SetBlueValue(lights.Value(outputColor.Blue), time.Duration(0))
		this.SetWhiteValue(lights.Value(outputColor.White), time.Duration(0))
		if this.IsAmberAvailable() {
			this.SetAmberValue(lights.Value(outputColor.Amber), time.Duration(0))
		}
		if this.IsUvAvailable() {
			this.SetUvValue(lights.Value(outputColor.Uv), time.Duration(0))
		}
	case *channellight.FourChannelLight:
		outputColor := v.GetOutputColor()

		this.SetFaderValue(lights.Value(outputColor.Brightness), time.Duration(0))
		this.SetRedValue(lights.Value(outputColor.Red), time.Duration(0))
		this.SetGreenValue(lights.Value(outputColor.Green), time.Duration(0))
		this.SetBlueValue(lights.Value(outputColor.Blue), time.Duration(0))
	}
}

func (this *baseRGBFixture) doGetValue(description string, channel lights.ChannelNo) lights.Value {
	if channel == channelNotSpecified {
		log.Infof("(%s) Requested '%s' channel when not implemented", this.accessory.Info.Name, description)
		return 0
	}
	return store.GetValue(channel)
}

func (this *baseRGBFixture) IsFaderAvailable() bool {
	return this.colorFixtureChannels.fader != channelNotSpecified
}

func (this *baseRGBFixture) IsWhiteAvailable() bool {
	return this.colorFixtureChannels.white != channelNotSpecified
}

func (this *baseRGBFixture) IsAmberAvailable() bool {
	return this.colorFixtureChannels.amber != channelNotSpecified
}

func (this *baseRGBFixture) IsUvAvailable() bool {
	return this.colorFixtureChannels.uv != channelNotSpecified
}

func (this *baseRGBFixture) GetFaderValue() lights.Value {
	return this.doGetValue("fader", this.colorFixtureChannels.fader)
}

func (this *baseRGBFixture) GetRedValue() lights.Value {
	return this.doGetValue("red", this.colorFixtureChannels.red)
}

func (this *baseRGBFixture) GetGreenValue() lights.Value {
	return this.doGetValue("green", this.colorFixtureChannels.green)
}

func (this *baseRGBFixture) GetBlueValue() lights.Value {
	return this.doGetValue("blue", this.colorFixtureChannels.blue)
}

func (this *baseRGBFixture) GetWhiteValue() lights.Value {
	return this.doGetValue("white", this.colorFixtureChannels.white)
}

func (this *baseRGBFixture) GetAmberValue() lights.Value {
	return this.doGetValue("amber", this.colorFixtureChannels.amber)
}

func (this *baseRGBFixture) GetUvValue() lights.Value {
	return this.doGetValue("uv", this.colorFixtureChannels.uv)
}

func (this *baseRGBFixture) doSetValue(description string, channel lights.ChannelNo, value lights.Value, fade time.Duration) {
	if channel == channelNotSpecified {
		log.Errorf("(%s) Attempted to set '%s' channel when not implemented", this.accessory.Info.Name, description)
		return
	}

	channelUpdater.GetChannelUpdater(channel).UpdateValueWithFade(this.doGetValue(description, channel), value, fade)
}

func (this *baseRGBFixture) SetFaderValue(value lights.Value, fade time.Duration) {
	this.doSetValue("fader", this.colorFixtureChannels.fader, value, fade)
}

func (this *baseRGBFixture) SetRedValue(value lights.Value, fade time.Duration) {
	this.doSetValue("red", this.colorFixtureChannels.red, value, fade)
}

func (this *baseRGBFixture) SetGreenValue(value lights.Value, fade time.Duration) {
	this.doSetValue("green", this.colorFixtureChannels.green, value, fade)
}

func (this *baseRGBFixture) SetBlueValue(value lights.Value, fade time.Duration) {
	this.doSetValue("blue", this.colorFixtureChannels.blue, value, fade)
}

func (this *baseRGBFixture) SetWhiteValue(value lights.Value, fade time.Duration) {
	this.doSetValue("white", this.colorFixtureChannels.white, value, fade)
}

func (this *baseRGBFixture) SetAmberValue(value lights.Value, fade time.Duration) {
	this.doSetValue("amber", this.colorFixtureChannels.amber, value, fade)
}

func (this *baseRGBFixture) SetUvValue(value lights.Value, fade time.Duration) {
	this.doSetValue("uv", this.colorFixtureChannels.uv, value, fade)
}

func (this *baseRGBFixture) SetColor(red, green, blue lights.Value, fade time.Duration) {
	this.SetRedValue(red, fade)
	this.SetGreenValue(green, fade)
	this.SetBlueValue(blue, fade)
}

func (this *baseRGBFixture) GetHomeKitAccessory() *hkAccessory.Accessory {
	return this.accessory
}

func (this *baseRGBFixture) ValueChange(change store.ValuesChange) {
	switch change.Channel {
	case this.colorFixtureChannels.fader, this.colorFixtureChannels.red, this.colorFixtureChannels.green, this.colorFixtureChannels.blue, this.colorFixtureChannels.white, this.colorFixtureChannels.amber, this.colorFixtureChannels.uv:
		log.Infof("(%s) Colour channel %d changed to %d (HomeKit update would happen here)", this.accessory.Info.Name.Value, change.Channel, change.Value)
	}
}
