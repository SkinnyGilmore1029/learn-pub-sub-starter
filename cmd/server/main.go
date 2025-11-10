package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Step 1: Declare connection string
	connStr := "amqp://guest:guest@localhost:5672/"

	// Step 2: Connect to RabbitMQ
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Step 3: Create a channel
	pauseChannel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer pauseChannel.Close()

	fmt.Println("✅ Successfully connected to RabbitMQ!")

	// 🟢 Step 4: Publish pause message to the exchange
	state := routing.PlayingState{IsPaused: true}

	err = pubsub.PublishJSON(
		pauseChannel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		state,
	)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	fmt.Println("📨 Published pause message to exchange:", routing.ExchangePerilDirect)

	// Step 5: Wait for Ctrl+C (SIGINT)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan // Block until Ctrl+C pressed

	fmt.Println("🛑 Shutting down server...")
}
