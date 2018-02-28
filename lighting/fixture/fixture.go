// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package fixture

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brutella/hc/accessory"
	"github.com/op/go-logging"
	"io/ioutil"
	"lighting/fixture/fixtureImpl"
	"lighting/fixtureType"
	"lighting/lights"
	"path"
)

var loadedFixtures = make(map[string]fixtureImpl.FixtureImpl)
var fixtureList []fixtureImpl.FixtureImpl

var log = logging.MustGetLogger("fixture")

const fixtureLocation = "/Users/chris/Development/Personal/lighting4/src/static/lighting/data/fixtures/"

type baseFixture struct {
	fixture      fixtureImpl.FixtureImpl
	Type         *fixtureType.FixtureType
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	FirstChannel lights.Address           `json:"firstChannel"`
	accessory    *accessory.Accessory
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

func (this *baseFixture) SetHomeKitAccessory(accessory *accessory.Accessory) {
	this.accessory = accessory
}

func (this *baseFixture) GetHomeKitAccessory() *accessory.Accessory {
	return this.accessory
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

	if this.Name == "" {
		return errors.New("no 'name' specified")
	}
	if this.FirstChannel.ChannelNo == 0 {
		return errors.New("no 'firstChannel' specified")
	}
	if temp.Type == "" {
		return errors.New("no 'type' specified")
	}

	typeOfFixture, err := fixtureType.GetFixtureType(temp.Type)
	if err != nil {
		return err
	}

	this.Type = typeOfFixture

	return nil
}

func GetFixtures() ([]fixtureImpl.FixtureImpl, error) {
	if len(fixtureList) > 0 {
		return fixtureList, nil
	}

	files, err := ioutil.ReadDir(fixtureLocation)
	if err != nil {
		return []fixtureImpl.FixtureImpl{}, err
	}

	var foundFixtures []fixtureImpl.FixtureImpl

	for _, fixtureFile := range files {
		fixture, err := doLoadFixture(fixtureFile.Name())
		if err != nil {
			log.Errorf("Could not load '%s' (%s)", fixtureFile.Name(), err)
		} else {
			foundFixtures = append(foundFixtures, fixture)
		}
	}

	fixtureList = foundFixtures

	return foundFixtures, nil
}

func doLoadFixture(fileName string) (fixtureImpl.FixtureImpl, error) {
	fixture, ok := loadedFixtures[fileName]
	if ok {
		return fixture, nil
	}

	raw, err := ioutil.ReadFile(path.Join(fixtureLocation, fileName))
	if err != nil {
		return nil, err
	}

	myBaseFixture := &baseFixture{}
	err = json.Unmarshal(raw, myBaseFixture)
	if err != nil {
		return nil, err
	}

	if fmt.Sprintf("%s.json", myBaseFixture.Name) != fileName {
		return nil, errors.New(fmt.Sprintf("fixture name (%s) and file name (%s) mismatch", myBaseFixture.Name, fileName))
	}


	switch myBaseFixture.Type.TypeKey {
	case "hex":
		fixture = fixtureImpl.NewChauvetHex(myBaseFixture)
	default:
		return nil, errors.New(fmt.Sprintf("no type implementation for '%s'", myBaseFixture.Type.TypeKey))
	}

	loadedFixtures[fileName] = fixture

	return fixture, nil
}

func GetFixture(fixtureKey string) (fixtureImpl.FixtureImpl, error) {
	fixture, err := doLoadFixture(fmt.Sprintf("%s.json", fixtureKey))
	if err != nil {
		return nil, err
	}

	return fixture, nil
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
