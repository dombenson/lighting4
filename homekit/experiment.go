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

	fivechLight := channellight.FiveChannelLight{}
	fivechLight.SetColor(lightbulb)

	sixchLight := channellight.SixChannelLight{}
	sixchLight.SetColor(lightbulb)

	sevenChColor := light.GetOutputColor()
	sixChColor := sixchLight.GetOutputColor()
	fiveChColor := fivechLight.GetOutputColor()
	fourChColor := fourchLight.GetOutputColor()

	inColor := light.GetColor()

	log.Println("------------------------------------------------------------")
	log.Println("HSL Input")
	log.Println("------------------------------------------------------------")
	log.Printf("H: %.0f, S: %.0f, B: %d\n", inColor.Hue, inColor.Saturation, inColor.Brightness)

	log.Println("------------------------------------------------------------")
	log.Println("7 Channel Light")
	log.Println("------------------------------------------------------------")
	log.Printf("R: %d, G: %d, B: %d, W: %d, A: %d, u: %d, b: %d\n", sevenChColor.Red, sevenChColor.Green, sevenChColor.Blue, sevenChColor.White, sevenChColor.Amber, sevenChColor.Uv, sevenChColor.Brightness)
	log.Println("------------------------------------------------------------")
	log.Println("6 Channel Light")
	log.Println("------------------------------------------------------------")
	log.Printf("R: %d, G: %d, B: %d, W: %d, A:%d, b: %d\n", sixChColor.Red, sixChColor.Green, sixChColor.Blue, sixChColor.White, sixChColor.Amber, sixChColor.Brightness)
	log.Println("------------------------------------------------------------")
	log.Println("5 Channel Light")
	log.Println("------------------------------------------------------------")
	log.Printf("R: %d, G: %d, B: %d, W: %d, b: %d\n", fiveChColor.Red, fiveChColor.Green, fiveChColor.Blue, fiveChColor.White, fiveChColor.Brightness)
	log.Println("------------------------------------------------------------")
	log.Println("4 Channel Light")
	log.Println("------------------------------------------------------------")
	log.Printf("R: %d, G: %d, B: %d, b: %d\n", fourChColor.Red, fourChColor.Green, fourChColor.Blue, fourChColor.Brightness)
	log.Println("------------------------------------------------------------")
}
