// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package routes

import (
	"encoding/json"
	"lighting/fixture"
	"net/http"
)

func FixturesList(w http.ResponseWriter, r *http.Request) {
	fixtures, err := fixture.GetFixtures()
	if err != nil {
		panic(err)
	}

	fixtureList := struct {
		Fixtures map[int][]string `json:"fixtures"`
	}{
		Fixtures: make(map[int][]string),
	}

	for _, f := range fixtures {
		fixtureList.Fixtures[f.GetFirstChannel().Universe] = append(fixtureList.Fixtures[f.GetFirstChannel().Universe], f.GetName())
	}

	jsonBytes, err := json.Marshal(fixtureList)
	if err != nil {
		panic(err)
	}

	w.Write(jsonBytes)
}
