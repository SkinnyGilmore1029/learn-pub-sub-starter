package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	fmt.Println("✅ Connected to RabbitMQ!")
	gamelogic.PrintServerHelp()

	// Declare the exchange to make sure it exists
	err = ch.ExchangeDeclare(
		routing.ExchangePerilDirect, // name
		"direct",                    // type
		true,                        // durable
		false,                       // autoDelete
		false,                       // internal
		false,                       // noWait
		nil,                         // args
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		"game_logs",
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to game logs: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	/* Publish a pause message
	state := routing.PlayingState{IsPaused: true}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, state)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	fmt.Println("📨 Published pause message!")
	Leaving this for references and notes
	*/
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			fmt.Println("📨 Sending pause message!")
			state := routing.PlayingState{IsPaused: true}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, state)
			if err != nil {
				log.Fatalf("Failed to publish message: %v", err)
			}

		case "resume":
			fmt.Println("📨 Sending resume message!")

			state := routing.PlayingState{IsPaused: false}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, state)
			if err != nil {
				log.Fatalf("Failed to publish message: %v", err)
			}

		case "quit":
			fmt.Println("🛑 Shutting down server...")
			return

		default:
			fmt.Println("Unknown Command")
		}

		/* Wait for Ctrl+C
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		fmt.Println("🛑 Shutting down server...")
		Leaving this for references and notes
		*/
	}
}
