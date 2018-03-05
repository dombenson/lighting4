// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"github.com/op/go-logging"
	"goji.io"
	"goji.io/pat"
	"lighting/amqp"
	"lighting/amqp/lightingControl"
	"lighting/amqp/lightingUpdates"
	"lighting/amqp/trackControl"
	"lighting/amqp/trackUpdates"
	"lighting/fixture"
	"lighting/homekit"
	"lighting/routes"
	"lighting/sequencer"
	"lighting/sequencer/staticSequences"
	"lighting/socket"
	"lighting/store"
	"lighting/tracks"
	"net/http"
	"os"
)

func main() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled)

	err := amqp.Init()
	if err != nil {
		panic(err)
	}
	defer amqp.Close()

	err = lightingControl.Start()
	if err != nil {
		panic(err)
	}

	err = lightingUpdates.Start()
	if err != nil {
		panic(err)
	}

	err = trackControl.Start()
	if err != nil {
		panic(err)
	}

	err = trackUpdates.Start()
	if err != nil {
		panic(err)
	}

	err = store.Sync()
	if err != nil {
		panic(err)
	}

	err = tracks.Sync()
	if err != nil {
		panic(err)
	}

	loadedFixtures, err := fixture.GetFixtures()
	if err != nil {
		panic(err)
	}
	homekit.RegisterFixtures(loadedFixtures)

	homekit.Start("DMX-Lights", "06592309")
	defer homekit.Stop()

	defer sequencer.Stop()

	staticSequences.Register(staticSequences.NewTwoColorSequence("Red and Blue",
	                         staticSequences.NewColor(255, 0, 0),
	                         staticSequences.NewColor(0, 0, 255)))
	staticSequences.Register(staticSequences.NewTwoColorSequence("Green and Blue",
	                         staticSequences.NewColor(0, 255, 0),
	                         staticSequences.NewColor(0, 0, 255)))
	staticSequences.Register(staticSequences.NewTwoColorSequence("Red and Green",
	                         staticSequences.NewColor(255, 0, 0),
	                         staticSequences.NewColor(0, 255, 0)))
	staticSequences.Register(staticSequences.NewTwoColorSequence("Mint and blackberry",
	                         staticSequences.NewColor(135, 14, 68),
	                         staticSequences.NewColor(128, 211, 155)))
	staticSequences.Register(staticSequences.NewTwoColorSequence("Lemon Blueberry",
	                         staticSequences.NewColor(237, 246, 25),
	                         staticSequences.NewColor(160, 70, 224)))
	staticSequences.Register(staticSequences.NewTwoColorSequence("UV!",
	                         staticSequences.NewUVColor(),
	                         staticSequences.NewUVColor()))

	staticSequences.Start()
	defer staticSequences.Stop()

	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/lighting/socket"), socket.Handler)

	mux.HandleFunc(pat.Get("/lighting/fixtures/list"), routes.FixturesList)

	staticFilesLocation := "/Users/chris/Development/Personal/lighting4/src/static"
	mux.Handle(pat.Get("/lighting/*"), http.FileServer(http.Dir(staticFilesLocation)))

	err = http.ListenAndServe("localhost:8000", mux)
	if err != nil {
		panic(err)
	}
}
