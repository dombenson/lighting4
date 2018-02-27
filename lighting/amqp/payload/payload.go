// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package payload

type Payload struct {
	Event string      `json:"event"`
}

type LightingPayload struct {
	Payload
	Universe int `json:"universe"`
}

func NewLightingPayload(event string, universe int) LightingPayload {
	if universe == 0 {
		universe = 1
	}

	return LightingPayload{
		Payload: Payload{event},
		Universe: universe,
	}
}
