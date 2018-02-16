// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package lightingUpdates

import (
	"encoding/json"
	"github.com/op/go-logging"
	amqpLib "github.com/streadway/amqp"
	"lighting/amqp"
	"lighting/amqp/payload"
	"lighting/lights"
	"lighting/store"
)

var log = logging.MustGetLogger("lighting.updates")

var updateQueue amqpLib.Queue

var started bool

type valueSetData struct {
	Channel lights.ChannelNo `json:"c"`
	Value   lights.Value     `json:"v"`
	SeqNo   int              `json:"s"`
}

type valueSetPayload struct {
	payload.Payload
	Data []valueSetData `json:"data"`
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

	err = channel.QueueBind(updateQueue.Name, "", "lighting.updates", false, nil)
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
				case "vs", "value-set", "vr", "value-requested":
					for _, v := range details.Data {
						if v.SeqNo > store.GetLastSeenHardwareSeqNo(v.Channel) {
							log.Debugf("Set %d to %d (%d)", v.Channel, v.Value, v.SeqNo)
							store.SetLastSeenHardwareSeqNo(v.Channel, v.SeqNo)
							store.SetValue(v.Channel, v.Value)
						}
					}
				case "hr", "hardware-reset":
					log.Info("Hardware Reset")
					store.Reset()
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
