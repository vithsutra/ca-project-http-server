package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type rabbitmqConnection struct {
	conn  *amqp.Connection
	chann *amqp.Channel
}

func NewRabbitmqConnection() *rabbitmqConnection {
	rabbitmqUrl := os.Getenv("RABBITMQ_URL")

	if rabbitmqUrl == "" {
		log.Fatalln("missing RABBITMQ_URL env variable")
	}
	conn, err := amqp.Dial(rabbitmqUrl)

	if err != nil {
		log.Fatalln("failed to connect to rabbitmq broker, Error:", err.Error())
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalln("failed to open rabbitmq channel, Error: ", err.Error())
	}

	log.Println("connected to the rabbitmq")

	return &rabbitmqConnection{
		conn,
		ch,
	}
}
