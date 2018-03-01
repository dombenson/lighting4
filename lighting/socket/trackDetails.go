// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"lighting/tracks"
)

type trackDetailsPayload struct {
	socketPayload
	Data trackDetailsValue `json:"data"`
}

type trackDetailsValue struct {
	IsPlaying bool   `json:"isPlaying"`
	Artist    string `json:"artist,omitempty"`
	Name      string `json:"name,omitempty"`
}

func (this *socketConnection) doSendTrackDetails(event string, trackState *tracks.TrackState) error {
	trackStateValue := trackDetailsValue{
		IsPlaying: false,
	}

	if trackState != nil {
		trackStateValue.IsPlaying = trackState.PlayerState == tracks.PlayerStatePlaying
		trackStateValue.Artist = trackState.Artist
		trackStateValue.Name = trackState.Name
	}

	this.mu.Lock()
	defer this.mu.Unlock()

	details := trackDetailsPayload{
		socketPayload: socketPayload {event},
		Data:          trackStateValue,
	}

	message, err := json.Marshal(details)
	if err != nil {
		return err
	}

	err = this.c.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil
}

func (this *socketConnection) processTrackDetails() error {
	currentTrackState := tracks.GetValue()

	err := this.doSendTrackDetails(trackDetails, currentTrackState)
	if err != nil {
		return err
	}

	log.Infof("(%d) 'trackDetails' sent", this.id)

	return err
}

func (this *socketConnection) notifyTrackChanged() tracks.ValueChangeCallback {
	return func(change tracks.ValuesChange) {
		err := this.doSendTrackDetails(notifyTrackChange, change.TrackState)
		if err != nil {
			log.Errorf("(%d) Error in track change notify (%s)", this.id, err)
		}

		return
	}
}
