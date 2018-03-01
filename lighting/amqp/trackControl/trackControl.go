// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package trackControl

import (
	"encoding/json"
	"github.com/op/go-logging"
	amqpLib "github.com/streadway/amqp"
	"lighting/amqp"
	"lighting/amqp/payload"
)

const controlExchange = "lighting.trackControl"

var log = logging.MustGetLogger("lighting.trackControl")

var started bool

func Start() error {
	if started {
		return nil
	}

	channel := amqp.GetChannel()

	err := channel.ExchangeDeclare(
		controlExchange,
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

	log.Info("AMQP started")

	started = true
	return nil
}

func RequestState() error {
	if !started {
		panic("lighting.trackControl exchange not started")
	}

	channel := amqp.GetChannel()

	requestData := payload.Payload{"rs"}

	jsonBytes, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	err = channel.Publish(controlExchange, "", false, false, amqpLib.Publishing {
		ContentType: "application/json",
		Body:        jsonBytes,
	})
	if err != nil {
		return err
	}

	return nil
}
