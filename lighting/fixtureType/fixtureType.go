// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixtureType

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lighting/lights"
	"path"
)

var loadedFixtureTypes map[string]*FixtureType

const fixtureTypeLocation = "/Users/chris/Development/Personal/lighting4/src/static/lighting/data/fixtureTypes/"

func init() {
	loadedFixtureTypes = make(map[string]*FixtureType)
}

type FixtureType struct {
	TypeKey                string             `json:"type"`
	ChannelCount           lights.ChannelNo   `json:"channelCount"`
	BlackoutChannelOffsets []lights.ChannelNo `json:"blackoutChannelOffsets"`
}

type FixtureProperty interface {
	GetOffsets() map[string]lights.ChannelNo
}

func GetFixtureType(typeKey string) (*FixtureType, error) {
	fixtureType, ok := loadedFixtureTypes[typeKey]
	if ok {
		return fixtureType, nil
	}

	raw, err := ioutil.ReadFile(path.Join(fixtureTypeLocation, fmt.Sprintf("%s.json", typeKey)))
	if err != nil {
		return nil, err
	}

	fixtureType = nil
	err = json.Unmarshal(raw, &fixtureType)
	if err != nil {
		return nil, err
	}

	return fixtureType, nil
}
