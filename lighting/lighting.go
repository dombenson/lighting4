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
	"lighting/sequencer"
	"lighting/socket"
	"lighting/store"
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

	err = store.Sync()
	if err != nil {
		panic(err)
	}

	err = sequencer.Start()
	if err != nil {
		panic(err)
	}
	defer sequencer.Stop()

	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/lighting/socket"), socket.Handler)

	staticFilesLocation := "/Users/chris/Development/Personal/lighting4/src/static"
	mux.Handle(pat.Get("/lighting/*"), http.FileServer(http.Dir(staticFilesLocation)))

	err = http.ListenAndServe("localhost:8000", mux)
	if err != nil {
		panic(err)
	}
}
