// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package sequencer

type Sequence interface {
	GetName() string
	Render()
	Stop()
}

var currentSequence Sequence

func SetSequence(sequence Sequence) {
	currentSequence = sequence
	currentSequence.Render()
}

func Stop() {
	if currentSequence != nil {
		currentSequence.Stop()
	}
}
