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
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Prompt user for a username
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Bad username: %v", err)
	}

	// Build queue name like "pause.username"
	queueName := routing.PauseKey + "." + username

	// Declare and bind queue
	ch, q, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect, // exchange
		queueName,                   // queue name
		routing.PauseKey,            // routing key
		pubsub.SimpleQueueTransient, // queue type
	)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}
	defer ch.Close()

	fmt.Printf("Client connected. Queue: %s (Exchange: %s)\n", q.Name, routing.ExchangePerilDirect)
	fmt.Println("Starting Peril client...")
}
