// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package amqp

import (
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

var log = logging.MustGetLogger("amqp")

func Init() error {
	if channel != nil || conn != nil {
		return nil
	}

	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}

	channel, err = conn.Channel()
	if err != nil {
		conn.Close()
		conn = nil
		return err
	}

	log.Info("Initialised")

	return nil
}

func GetChannel() *amqp.Channel {
	if channel == nil {
		panic("attempted to get channel when not initialised")
	}
	return channel
}

func Close() {
	if channel != nil {
		log.Info("Closing AMQP channel")
		err := channel.Close()
		if err != nil {
			log.Error(err)
		}
		channel = nil
	}

	if conn != nil {
		log.Info("Closing AMQP connection")
		err := conn.Close()
		if err != nil {
			log.Error(err)
		}

		conn = nil
	}

	log.Info("Finished closing AMQP")
}