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
		FixtureImpl:      fixture,
		baseRGBFixture:   &baseRGBFixture{
			redChannel:   fixture.GetFirstChannel() + 1,
			greenChannel: fixture.GetFirstChannel() + 2,
			blueChannel:  fixture.GetFirstChannel() + 3,
		},
	}

	return &rgbFixture
}

