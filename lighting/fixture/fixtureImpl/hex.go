// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

import "lighting/lights"

type ChauvetHex struct {
	FixtureImpl
	*baseRGBFixture
}

func NewChauvetHex(fixture FixtureImpl) *ChauvetHex {
	faderChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo)
	redChannel   := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 1)
	greenChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 2)
	blueChannel  := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 3)
	amberChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 4)
	whiteChannel := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 5)
	uvChannel    := lights.NewAddress(fixture.GetFirstChannel().Universe, fixture.GetFirstChannel().ChannelNo + 6)

	rgbFixture := ChauvetHex {
		FixtureImpl:    fixture,
		baseRGBFixture: newBaseRGBFixture(fixture, colorFixtureChannels{
			fader: &faderChannel,
			red:   &redChannel,
			green: &greenChannel,
			blue:  &blueChannel,
			amber: &amberChannel,
			white: &whiteChannel,
			uv:    &uvChannel,
		}),
	}

	return &rgbFixture
}

