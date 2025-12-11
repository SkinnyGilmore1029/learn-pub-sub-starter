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
	fmt.Println("Starting Peril client...")
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}

	newstate := gamelogic.NewGameState(username)

	// Subscribe to pause/resume messages
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,   // exchange
		routing.PauseKey+"."+username, // queue name (unique for each client)
		routing.PauseKey,              // routing key (pause)
		pubsub.SimpleQueueTransient,   // queue type
		handlerPause(newstate),        // handler function
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause messages: %v", err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			if err := newstate.CommandSpawn(input); err != nil {
				fmt.Println("Error:", err)
			}
		case "move":
			_, err := newstate.CommandMove(input)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "status":
			newstate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			fmt.Println("🛑 Shutting down client...")
			return
		default:
			fmt.Println("Unknown Command")
			continue
		}
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(state routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(state)
	}
}
