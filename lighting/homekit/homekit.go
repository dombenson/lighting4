// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package homekit

import (
	"github.com/brutella/hc"
	hkAccessory "github.com/brutella/hc/accessory"
	"github.com/op/go-logging"
	"lighting/fixture/fixtureImpl"
)

var log = logging.MustGetLogger("homekit")

var homeKitTransport hc.Transport
var initialised = false

var fixtureMap = make(map[string]*hkAccessory.Accessory)
var fixtures []*hkAccessory.Accessory

func RegisterFixtures(fixturesToRegister []fixtureImpl.FixtureImpl) {
	if initialised {
		panic("Attempted to register fixture after HomeKit started")
	}

	for _, fixtureToRegister := range fixturesToRegister {
		accessory := fixtureToRegister.GetHomeKitAccessory()
		if accessory != nil {

			_, exists := fixtureMap[fixtureToRegister.GetName()]
			if exists {
				panic("Attempted to register a fixture that has already been registered")
			}

			fixtures = append(fixtures, accessory)
			fixtureMap[fixtureToRegister.GetName()] = accessory
		}
	}
}

func Start(name, pin string) error {
	bridgeInfo := hkAccessory.Info{
		Name: name,
	}
	bridge := hkAccessory.NewBridge(bridgeInfo)

	config := hc.Config{
		Pin: pin,
	}

	t, err := hc.NewIPTransport(config, bridge.Accessory, fixtures...)
	if err != nil {
		return err
	}

	homeKitTransport = t
	go homeKitTransport.Start()
	initialised = true

	return nil
}

func Stop() {
	if initialised {
		homeKitTransport.Stop()
	}
}
