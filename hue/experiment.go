// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main
import (
	"fmt"
	"github.com/cjcormack/gohue"
)
func main() {
	// It is recommended that you save the username from bridge.CreateUser
	// so you don't have to press the link button every time and re-auth.
	// When CreateUser is called it will print the generated user token.
	bridgesOnNetwork, err := hue.FindBridges()
	if err != nil {
		panic(err)
	}
	bridge := bridgesOnNetwork[0]
	//username, err := bridge.CreateUser("chrishueexperiment")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(username)

	err = bridge.Login("beYSuST8-9JxcFAhvls2dUj6TZpZXIYlnEILtq3C")
	if err != nil {
		panic(err)
	}

	groups, err := bridge.GetGroups()
	if err != nil {
		panic(err)
	}
	for _, group := range groups {
		//light.SetBrightness(100)
		//light.ColorLoop(true)
		fmt.Println(group.Name)
		for _, light := range group.Lights {
			fmt.Println(light.Index, light)
		}
	}

	//nightstandLight, _ := bridge.GetLightByName("Nightstand")
	//nightstandLight.Blink(5)
	//nightstandLight.SetName("Bedroom Lamp")

	//lights[0].SetColor(hue.RED)
	//lights[1].SetColor(hue.BLUE)
	//lights[2].SetColor(hue.GREEN)
	//
	//for _, light := range lights {
	//	light.Off()
	//}
}
