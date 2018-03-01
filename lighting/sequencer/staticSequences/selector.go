// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package staticSequences

import (
	"github.com/op/go-logging"
	"lighting/sequencer"
	"lighting/tracks"
	"math/rand"
	"time"
)

var log = logging.MustGetLogger("static-Sequencer")

var r = rand.New(rand.NewSource(time.Now().Unix()))

var unregisterFn func()

var currentSequenceIndex = -1

var sequences []sequencer.Sequence

func Start() {
	unregisterFn = tracks.Subscribe(TrackChange)
	ChooseRandomSequence()
}

func ChooseRandomSequence() {
	if len(sequences) > 0 {
		nextSequenceIndex := r.Intn(len(sequences))

		// attempt to find a different sequence (@todo improve this approach)
		for loopCount := 0; currentSequenceIndex == nextSequenceIndex && loopCount < 5; loopCount++ {
			nextSequenceIndex = r.Intn(len(sequences))
		}

		nextSequence := sequences[nextSequenceIndex]
		currentSequenceIndex = nextSequenceIndex

		sequencer.SetSequence(nextSequence)
	}
}

func Register(sequence sequencer.Sequence) {
	sequences = append(sequences, sequence)
}

func Stop()  {
	unregisterFn()
}

func TrackChange(change tracks.ValuesChange) {
	log.Info("Track changed")

	if !change.TrackChanged {
		return
	}

	ChooseRandomSequence()
}
