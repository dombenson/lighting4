// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
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

	hc.OnTermination(func(){
		<-t.Stop()
	})

	t.Start()
}

func convertColorsForLight(lightbulb *service.Lightbulb) {
	hue := lightbulb.Hue.GetValue()
	saturation := lightbulb.Saturation.GetValue()
	brightness := lightbulb.Brightness.GetValue()

	log.Printf("H: %.0f, S: %.0f, B: %d\n", hue, saturation, brightness)

	red := calculateRed(hue, saturation, brightness)
	green := calculateGreen(hue, saturation, brightness)
	blue := calculateBlue(hue, saturation, brightness)
	white := calculateWhite(hue, saturation, brightness)
	amber := calculateAmber(hue, saturation, brightness)
	uv := calculateUv(hue, saturation, brightness)

	log.Printf("R: %d, G: %d, B: %d, W: %d, A: %d, u: %d\n", red, green, blue, white, amber, uv)
	log.Println("------------------------------------------------------------")
}

func calculateRed(hue, saturation float64, brightness int) int {
	return 0
}

func calculateGreen(hue, saturation float64, brightness int) int {
	return 0
}

func calculateBlue(hue, saturation float64, brightness int) int {
	return 0
}

func calculateWhite(hue, saturation float64, brightness int) int {
	return 0
}

func calculateAmber(hue, saturation float64, brightness int) int {
	return 0
}

func calculateUv(hue, saturation float64, brightness int) int {
	return 0
}