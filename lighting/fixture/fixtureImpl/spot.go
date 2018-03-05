// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import "lighting/lights"

type Spot struct {
	FixtureImpl
	*baseRGBFixture
}

func NewSpot(fixture FixtureImpl) *Spot {
	redChannel   := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 1)
	greenChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 2)
	blueChannel  := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 3)

	rgbFixture := Spot {
		FixtureImpl:    fixture,
		baseRGBFixture: newBaseRGBFixture(fixture, colorFixtureChannels{
			red:   &redChannel,
			green: &greenChannel,
			blue:  &blueChannel,
		}),
	}

	return &rgbFixture
}

