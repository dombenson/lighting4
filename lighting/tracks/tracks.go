// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package tracks

import (
	"github.com/op/go-logging"
	"lighting/amqp/trackControl"
)

var log = logging.MustGetLogger("tracks")

const (
	PlayerStatePlaying = "Playing"
	PlayerStatePaused  = "Paused"
)

type TrackState struct {
	PlayerState string `json:"ps"`
	Name        string `json:"n"`
	Artist      string `json:"a"`
}

var currentState *TrackState

func Sync() error {
	err := trackControl.RequestState()
	if err != nil {
		return err
	}

	return nil
}

func SetValue(trackState *TrackState) {
	if trackState == nil {
		log.Info("State unset")
	} else {
		log.Infof("State changed to %s (%s, %s)", trackState.PlayerState, trackState.Name, trackState.Artist)
	}

	trackChanged := false

	if currentState == nil {
		if trackState != nil {
			trackChanged = true
		}
	} else {
		if currentState.Name != trackState.Name || currentState.Artist != trackState.Artist {
			trackChanged = true
		}
	}

	currentState = trackState

	notify(trackChanged, currentState)
}

func GetValue() *TrackState {
	return currentState
}