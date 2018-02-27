// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixture

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"lighting/fixture/fixtureImpl"
	"lighting/fixtureType"
	"lighting/lights"
	"path"
)

var loadedFixtures = make(map[string]fixtureImpl.FixtureImpl)

const fixtureLocation = "/Users/chris/Development/Personal/lighting4/src/static/lighting/data/fixtures/"

type baseFixture struct {
	fixture      fixtureImpl.FixtureImpl
	Type         *fixtureType.FixtureType
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	FirstChannel lights.Address           `json:"firstChannel"`
}

func (this *baseFixture) GetType() *fixtureType.FixtureType {
	return this.Type
}

func (this *baseFixture) GetName() string {
	return this.Name
}

func (this *baseFixture) GetDescription() string {
	return this.Description
}

func (this *baseFixture) GetFirstChannel() lights.Address {
	return this.FirstChannel
}

func (this *baseFixture) UnmarshalJSON(data []byte) error {
	temp := struct{
		Name         *string           `json:"name"`
		Description  *string           `json:"description"`
		FirstChannel *lights.ChannelNo `json:"firstChannel"`
		Universe     *int              `json:"universe"`
		Type         string            `json:"type"`
	}{
		Name:         &this.Name,
		Description:  &this.Description,
		FirstChannel: &this.FirstChannel.ChannelNo,
		Universe:     &this.FirstChannel.Universe,
	}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	if this.FirstChannel.Universe == 0 {
		this.FirstChannel.Universe = 1
	}

	typeOfFixture, err := fixtureType.GetFixtureType(temp.Type)
	if err != nil {
		return err
	}

	this.Type = typeOfFixture

	return nil
}

func GetFixture(fixtureKey string) (fixtureImpl.FixtureImpl, error) {
	fixture, ok := loadedFixtures[fixtureKey]
	if ok {
		return fixture, nil
	}

	raw, err := ioutil.ReadFile(path.Join(fixtureLocation, fmt.Sprintf("%s.json", fixtureKey)))
	if err != nil {
		return nil, err
	}

	myBaseFixture := &baseFixture{}
	err = json.Unmarshal(raw, myBaseFixture)
	if err != nil {
		return nil, err
	}

	switch myBaseFixture.Type.TypeKey {
	case "hex":
		return fixtureImpl.NewChauvetHex(myBaseFixture), nil
	default:
		return nil, errors.New(fmt.Sprintf("no type implementation for '%s'", myBaseFixture.Type.TypeKey))
	}
}

func GetRGBFixture(fixtureKey string) (fixtureImpl.RGBFixtureImpl, error) {
	fixture, err := GetFixture(fixtureKey)
	if err != nil {
		return nil, err
	}

	rgbFixture, ok := fixture.(fixtureImpl.RGBFixtureImpl)
	if !ok {
		return nil, errors.New(fmt.Sprintf("not an rgb fixture '%s'", fixtureKey))
	}

	return rgbFixture, nil
}
