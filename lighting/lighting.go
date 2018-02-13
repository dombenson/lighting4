// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"lighting/amqp"
	"lighting/amqp/lightingControl"
	"lighting/amqp/lightingUpdates"
	"lighting/store"
	"log"
)

func main() {
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

	err = lightingControl.SetValue(1, 235)
	if err != nil {
		panic(err)
	}

	observer := make(chan store.ValuesChange)

	store.Attach(observer)

	for {
		select {
		case value := <-observer:
			log.Printf("Observed channel %d changing to %d", value.Channel, value.Value)
			if value.Value == 255 {
				return
			}

			err = lightingControl.SetValue(1, 255)
			if err != nil {
				panic(err)
			}
		}
	}

	//mux := goji.NewMux()
	//
	//err = http.ListenAndServe("localhost:8000", mux)
	//if err != nil {
	//	panic(err)
	//}
}
