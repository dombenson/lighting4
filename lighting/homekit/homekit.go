// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("homekit")

var homeKitTransport hc.Transport
var initialised = false

var fixtureMap = make(map[string]*accessory.Accessory)
var fixtures []*accessory.Accessory

func RegisterFixture(lightKey string, fixture *accessory.Accessory) {
	if initialised {
		panic("Attempted to register fixture after HomeKit started")
	}

	_, exists := fixtureMap[lightKey]
	if exists {
		panic("Attempted to register a fixture that has already been registered")
	}

	fixtures = append(fixtures, fixture)
	fixtureMap[lightKey] = fixture
}

func Start(name, pin string) error {
	bridgeInfo := accessory.Info{
		Name: name,
	}
	bridge := accessory.NewBridge(bridgeInfo)

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
