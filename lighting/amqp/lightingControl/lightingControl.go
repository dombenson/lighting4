// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package lightingControl

import (
	"encoding/json"
	"github.com/op/go-logging"
	amqpLib "github.com/streadway/amqp"
	"lighting/amqp"
	"lighting/amqp/payload"
	"lighting/lights"
)

const controlExchange = "lighting.control"

var log = logging.MustGetLogger("lighting.control")

var started bool

type setValueData struct {
	Channel lights.ChannelNo `json:"c"`
	Value   lights.Value     `json:"v"`
}

type setValuePayload struct {
	payload.LightingPayload
	Data setValueData `json:"data"`
}

type requestValuesData []lights.ChannelNo

type requestValuesPayload struct {
	payload.LightingPayload
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

	log.Info("AMQP started")

	started = true
	return nil
}

func SetValue(address lights.Address, value lights.Value) error {
	if !started {
		panic("lighting.control exchange not started")
	}

	channel := amqp.GetChannel()

	setData := setValuePayload{
		LightingPayload: payload.NewLightingPayload("sv", address.Universe),
		Data: setValueData {
			Channel: address.ChannelNo,
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

func RequestValues(universe int) error {
	if !started {
		panic("lighting.control exchange not started")
	}

	channel := amqp.GetChannel()

	requestData := requestValuesPayload{
		LightingPayload: payload.NewLightingPayload("rv", universe),
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