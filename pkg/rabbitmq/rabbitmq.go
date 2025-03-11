package rabbitmq

import (
	"errors"
	"os"

	"github.com/streadway/amqp"
)

type rabbitmqRepo struct {
	conn  *amqp.Connection
	chann *amqp.Channel
}

func NewRabbitmqRepo(conn *amqp.Connection, chann *amqp.Channel) *rabbitmqRepo {
	return &rabbitmqRepo{
		conn,
		chann,
	}
}

func (repo *rabbitmqRepo) SendEmail(data []byte) error {
	queueName := os.Getenv("QUEUE_NAME")

	if queueName == "" {
		return errors.New("missing QUEUE_NAME env variable")
	}

	err := repo.chann.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
			Priority:    10,
		},
	)

	return err
}
