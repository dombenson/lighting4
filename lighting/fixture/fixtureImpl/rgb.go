// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureImpl

type GenericRGB struct {
	FixtureImpl
	*baseRGBFixture
}

func NewGenericRGB(fixture FixtureImpl) *GenericRGB {
	rgbFixture := GenericRGB {
		FixtureImpl:      fixture,
		baseRGBFixture:   &baseRGBFixture{
			redChannel:   fixture.GetFirstChannel() + 0,
			greenChannel: fixture.GetFirstChannel() + 1,
			blueChannel:  fixture.GetFirstChannel() + 2,
		},
	}

	return &rgbFixture
}

