// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package trackUpdates

import (
	"encoding/json"
	"github.com/op/go-logging"
	amqpLib "github.com/streadway/amqp"
	"lighting/amqp"
	"lighting/amqp/payload"
	"lighting/tracks"
)

var log = logging.MustGetLogger("lighting.trackUpdates")

var updateQueue amqpLib.Queue

var started bool

type playerStateData struct {
	PlayerState string `json:"ps"`
	Name        string `json:"n"`
	Artist      string `json:"a"`
}

type valueSetPayload struct {
	payload.Payload
	Data *playerStateData `json:"data,omitempty"`
}

func Start() error {
	if started {
		return nil
	}

	channel := amqp.GetChannel()

	err := channel.ExchangeDeclare(
		"lighting.updates",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	updateQueue, err = channel.QueueDeclare(
		"", // name
		true,   // durable
		true,   // delete when unused
		true,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	err = channel.QueueBind(updateQueue.Name, "", "lighting.trackUpdates", false, nil)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		updateQueue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			details := &valueSetPayload{}

			err := json.Unmarshal(d.Body, &details)
			if err != nil {
				log.Error(err)
			} else {
				switch details.Event {
				case "ps", "playback-state", "sr", "state-requested", "tc", "track-changed":
					if details.Data == nil || details.Data.PlayerState == "" {
						tracks.SetValue(nil)
					} else {
						tracks.SetValue(&tracks.TrackState{
							PlayerState: details.Data.PlayerState,
							Name:        details.Data.Name,
							Artist:      details.Data.Artist,
						})
					}
				default:
					log.Error("Unsupported message", details)
				}
			}
		}
	}()

	started = true
	log.Info("AMQP started")

	return nil
}

