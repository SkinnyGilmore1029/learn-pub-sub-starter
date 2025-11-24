package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // enum: SimpleQueueDurable or SimpleQueueTransient
) (*amqp.Channel, amqp.Queue, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		queueType == SimpleQueueDurable,   // durable
		queueType == SimpleQueueTransient, // autoDelete
		queueType == SimpleQueueTransient, // exclusive
		false,                             // noWait
		nil,                               // args
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	if err := ch.QueueBind(
		q.Name,
		key,
		exchange,
		false, // noWait
		nil,   // args
	); err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, q, nil
}
