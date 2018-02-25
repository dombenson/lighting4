// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package sequencer

import (
	"lighting/fixture"
	"log"
	"time"
)

var exampleTicker *time.Ticker

func Start() error {
	var err error
	exampleTicker, err = exampleSequence()
	if err != nil {
		return err
	}

	return nil
}

func Stop() {
	exampleTicker.Stop()
}

func exampleSequence() (*time.Ticker, error) {
	rgb1, err := fixture.GetRGBFixture("rear-left")
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(10 * time.Second)

	go func() {
		tickCount := 0

		for range ticker.C {
			tickCount++
			switch tickCount {
			case 1:
				rgb1.SetColor(255, 0, 0, 1 * time.Second)
			case 2:
				log.Println()
				rgb1.SetColor(0, 255, 0, 1 * time.Second)
			case 3:
				log.Println()
				rgb1.SetColor(0, 0, 255, 1 * time.Second)
				tickCount = 0
			}
		}
	}()

	return ticker, nil
}
