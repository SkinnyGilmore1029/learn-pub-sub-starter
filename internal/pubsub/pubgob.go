package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

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

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		gobUnmarshaller[T],
	)
}

func gobUnmarshaller[T any](data []byte) (T, error) {
	var target T
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&target)
	return target, err
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
) error {
	ch, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("could not declare and bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	// **Use the generic unmarshaller passed in**
	go func() {
		defer ch.Close()
		defer fmt.Print("> ")
		for msg := range msgs {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("could not unmarshal message: %v\n", err)
				continue
			}
			switch handler(target) {
			case Ack:
				msg.Ack(false)
				fmt.Println("Ack")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("NackDiscard")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("NackRequeue")
			}
		}
	}()

	return nil
}

/*


		need a decoder from gob not json
		----------------------------------

		unmarshaller := func(data []byte) (T, error) {
			var target T
			err := json.Unmarshal(data, &target)
			return target, err
		}

go func() {
			defer ch.Close()
			defer fmt.Print("> ")
			for msg := range msgs {
				target, err := unmarshaller(msg.Body)
				if err != nil {
					fmt.Printf("could not unmarshal message: %v\n", err)
					continue
				}
				switch handler(target) {
				case Ack:
					msg.Ack(false)
					fmt.Println("Ack")
				case NackDiscard:
					msg.Nack(false, false)
					fmt.Println("NackDiscard")
				case NackRequeue:
					msg.Nack(false, true)
					fmt.Println("NackRequeue")
				}
			}
		}()



func encode(gl GameLog) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(gl); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decode(data []byte) (GameLog, error) {
	var gl GameLog

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&gl); err != nil {
		return GameLog{}, err
	}

	return gl, nil
}


*/
