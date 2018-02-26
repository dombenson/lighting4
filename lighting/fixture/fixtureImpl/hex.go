// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

type ChauvetHex struct {
	FixtureImpl
	*baseRGBFixture
}

func NewChauvetHex(fixture FixtureImpl) *ChauvetHex {
	rgbFixture := ChauvetHex {
		FixtureImpl:    fixture,
		baseRGBFixture: newBaseRGBFixture(fixture.GetName(), colorFixtureChannels{
			fader: fixture.GetFirstChannel(),
			red:   fixture.GetFirstChannel() + 1,
			green: fixture.GetFirstChannel() + 2,
			blue:  fixture.GetFirstChannel() + 3,
		}),
	}

	return &rgbFixture
}

