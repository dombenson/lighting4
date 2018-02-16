// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"channellight"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"log"
)

func main() {
	bridgeInfo := accessory.Info{
		Name: "My Bridge",
	}
	bridge := accessory.NewBridge(bridgeInfo)

	light1 := accessory.NewLightbulb(accessory.Info{
		Name: "Light One",
	})

	config := hc.Config{
		Pin: "00102003",
	}
	t, err := hc.NewIPTransport(config, bridge.Accessory, light1.Accessory)
	if err != nil {
		log.Panic(err)
	}

	light1.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			log.Println("Client changed switch to on")
		} else {
			log.Println("Client changed switch to off")
		}
	})

	light1.Lightbulb.Brightness.OnValueRemoteUpdate(func(b int) {
		convertColorsForLight(light1.Lightbulb)
	})

	light1.Lightbulb.Hue.OnValueRemoteUpdate(func(h float64) {
		convertColorsForLight(light1.Lightbulb)
	})

	light1.Lightbulb.Saturation.OnValueRemoteUpdate(func(s float64) {
		convertColorsForLight(light1.Lightbulb)
	})

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}

func convertColorsForLight(lightbulb *service.Lightbulb) {
	light := channellight.SevenChannelLight{}
	light.SetColor(lightbulb)

	fourchLight := channellight.FourChannelLight{}
	fourchLight.SetColor(lightbulb)

	log.Println("------------------------------------------------------------")
	log.Println("7 Channel Light")
	log.Println("------------------------------------------------------------")
	log.Printf("H: %.0f, S: %.0f, B: %d\n", light.GetColor().Hue, light.GetColor().Saturation, light.GetColor().Brightness)

	log.Printf("R: %d, G: %d, B: %d, W: %d, A: %d, u: %d, b: %d\n", light.GetOutputColor().Red, light.GetOutputColor().Green, light.GetOutputColor().Blue, light.GetOutputColor().White, light.GetOutputColor().Amber, light.GetOutputColor().Uv, light.GetOutputColor().Brightness)
	log.Println("------------------------------------------------------------")

	log.Println("------------------------------------------------------------")
	log.Println("4 Channel Light")
	log.Println("------------------------------------------------------------")
	log.Printf("H: %.0f, S: %.0f, B: %d\n", fourchLight.GetColor().Hue, fourchLight.GetColor().Saturation, fourchLight.GetColor().Brightness)

	log.Printf("R: %d, G: %d, B: %d, b: %d\n", fourchLight.GetOutputColor().Red, fourchLight.GetOutputColor().Green, fourchLight.GetOutputColor().Blue, fourchLight.GetOutputColor().Brightness)
	log.Println("------------------------------------------------------------")
}
