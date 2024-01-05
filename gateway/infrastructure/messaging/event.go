package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		// name
		"logs_topic",
		// kind
		"topic",
		// durable
		true,
		// auto-deleted
		false,
		// inernal
		false,
		// no-wait
		false,
		// arguments
		nil,
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		// name
		"",
		// durable
		false,
		// delete when unused
		false,
		// exclusive
		true,
		// no-wait
		false,
		// arguments
		nil,
	)
}
