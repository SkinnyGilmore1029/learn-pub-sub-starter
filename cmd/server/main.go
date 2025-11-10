package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	// Publish a pause message
	state := routing.PlayingState{IsPaused: true}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, state)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	fmt.Println("📨 Published pause message!")

	// Wait for Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("🛑 Shutting down server...")
}
