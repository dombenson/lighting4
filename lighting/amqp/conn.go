// Copyright 2018 Christopher Cormack. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package amqp

import (
	"github.com/streadway/amqp"
	"log"
)

var conn *amqp.Connection
var channel *amqp.Channel

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

	log.Println("AMQP initialised")

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
		log.Println("Closing AMQP channel")
		err := channel.Close()
		if err != nil {
			log.Println(err)
		}
		channel = nil
	}

	if conn != nil {
		log.Println("Closing AMQP connection")
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}

		conn = nil
	}

	log.Println("Finished closing AMQP")
}