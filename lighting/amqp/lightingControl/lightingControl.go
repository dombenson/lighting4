// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package lightingControl

import (
	"encoding/json"
	amqpLib "github.com/streadway/amqp"
	"lighting/amqp"
	"lighting/amqp/payload"
	"lighting/lights"
	"log"
)

const controlExchange = "lighting.control"

var started bool

type setValueData struct {
	Channel lights.ChannelNo `json:"c"`
	Value   lights.Value     `json:"v"`
}

type setValuePayload struct {
	payload.Payload
	Data setValueData `json:"data"`
}

type requestValuesData []lights.ChannelNo

type requestValuesPayload struct {
	payload.Payload
	Data requestValuesData `json:"data"`
}

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

	log.Println("Lighting Control AMQP started")

	started = true
	return nil
}

func SetValue(channelNo lights.ChannelNo, value lights.Value) error {
	if !started {
		panic("lighting.control exchange not started")
	}

	channel := amqp.GetChannel()

	setData := setValuePayload{
		Payload: payload.Payload{"sv"},
		Data: setValueData {
			Channel: channelNo,
			Value:   value,
		},
	}

	jsonBytes, err := json.Marshal(setData)
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

func RequestValues() error {
	if !started {
		panic("lighting.control exchange not started")
	}

	channel := amqp.GetChannel()

	requestData := requestValuesPayload{
		Payload: payload.Payload{"rv"},
	}

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