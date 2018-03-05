// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import "lighting/lights"

type Stairville struct {
	FixtureImpl
	*baseRGBFixture
}

func NewStairville(fixture FixtureImpl) *Stairville {
	faderChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo)
	redChannel   := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 1)
	greenChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 2)
	blueChannel  := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 3)

	rgbFixture := Stairville {
		FixtureImpl:    fixture,
		baseRGBFixture: newBaseRGBFixture(fixture, colorFixtureChannels{
			fader: &faderChannel,
			red:   &redChannel,
			green: &greenChannel,
			blue:  &blueChannel,
		}),
	}

	return &rgbFixture
}

