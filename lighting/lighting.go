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
	"lighting/channelUpdater"
	"lighting/socket"
	"lighting/store"
	"net/http"
	"os"
	"time"
)

func main() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.DEBUG, "")

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

	mux := goji.NewMux()

	staticFilesLocation := "/Users/chris/Development/Personal/lighting4/src/static"

	mux.HandleFunc(pat.Get("/lighting/socket"), socket.Handler)

	mux.Handle(pat.Get("/lighting/*"), http.FileServer(http.Dir(staticFilesLocation)))

	channelUpdater.GetChannelUpdater(1).UpdateValueWithFade(10, 222, 10 * time.Second)

	err = http.ListenAndServe("localhost:8000", mux)
	if err != nil {
		panic(err)
	}
}
