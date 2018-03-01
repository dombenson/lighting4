// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package staticSequences

import (
	"lighting/fixture"
	"lighting/fixture/fixtureImpl"
	"lighting/lights"
	"time"
)

type Color struct {
	Red   lights.Value
	Green lights.Value
	Blue  lights.Value
}

func NewColor(red, green, blue lights.Value) Color {
	return Color{
		Red:   red,
		Green: green,
		Blue:  blue,
	}
}

var colorOneFixtureKeys = []string{"tv-left"}
var colorTwoFixtureKeys = []string{"tv-right"}

type TwoColorSequence struct {
	Name             string
	ColorOne         Color
	ColorTwo         Color
	colorOneFixtures []fixtureImpl.RGBFixtureImpl
	colorTwoFixtures []fixtureImpl.RGBFixtureImpl
}

func (this *TwoColorSequence) GetName() string {
	return this.Name
}

func (this *TwoColorSequence) Render() {
	this.setColorOnFixtures(this.colorOneFixtures, this.ColorOne)
	this.setColorOnFixtures(this.colorTwoFixtures, this.ColorTwo)
}

func (this *TwoColorSequence) Stop() {
	this.setColorOnFixtures(this.colorOneFixtures, NewColor(0, 0, 0))
	this.setColorOnFixtures(this.colorTwoFixtures, NewColor(0, 0, 0))
}

func (this *TwoColorSequence) setColorOnFixtures(fixtures []fixtureImpl.RGBFixtureImpl, color Color) {
	for _, colorFixture := range fixtures {
		colorFixture.SetColor(color.Red, color.Green, color.Blue, 1 * time.Second)
	}
}

func NewTwoColorSequence(name string, colorOne, colorTwo Color) *TwoColorSequence {
	var colorOneFixtures []fixtureImpl.RGBFixtureImpl
	for _, colorOneFixtureKey := range colorOneFixtureKeys {
		colorOneFixture, err := fixture.GetRGBFixture(colorOneFixtureKey)
		if err != nil {
			panic(err)
		}

		colorOneFixtures = append(colorOneFixtures, colorOneFixture)
	}

	var colorTwoFixtures []fixtureImpl.RGBFixtureImpl
	for _, colorTwoFixtureKey := range colorTwoFixtureKeys {
		colorTwoFixture, err := fixture.GetRGBFixture(colorTwoFixtureKey)
		if err != nil {
			panic(err)
		}

		colorTwoFixtures = append(colorTwoFixtures, colorTwoFixture)
	}

	return &TwoColorSequence{
		Name:             name,
		ColorOne:         colorOne,
		ColorTwo:         colorTwo,
		colorOneFixtures: colorOneFixtures,
		colorTwoFixtures: colorTwoFixtures,
	}
}