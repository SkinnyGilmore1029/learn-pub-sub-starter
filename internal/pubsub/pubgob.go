package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// make a buffer to hold the gob data
	var buf bytes.Buffer

	// create a new encoder that will write to the buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(val); err != nil {
		return err
	}

	// return the publish
	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buf.Bytes(),
	})
}
