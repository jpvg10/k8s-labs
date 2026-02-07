package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type Todo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")

	nc, err := nats.Connect(natsURL)
	if err != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err)
		return
	}
	defer nc.Close()

	fmt.Printf("Connected to NATS at %s\n", natsURL)

	nc.Subscribe("todos.created", func(msg *nats.Msg) {
		var todo Todo
		if err := json.Unmarshal(msg.Data, &todo); err != nil {
			fmt.Printf("Error unmarshaling created todo: %v\n", err)
			return
		}
		fmt.Printf("[TODO CREATED] ID: %d, Description: %s, Done: %v\n", todo.Id, todo.Description, todo.Done)
	})

	nc.Subscribe("todos.updated", func(msg *nats.Msg) {
		var todo Todo
		if err := json.Unmarshal(msg.Data, &todo); err != nil {
			fmt.Printf("Error unmarshaling updated todo: %v\n", err)
			return
		}
		fmt.Printf("[TODO UPDATED] ID: %d, Description: %s, Done: %v\n", todo.Id, todo.Description, todo.Done)
	})

	fmt.Println("Listening for NATS messages on 'todos.created' and 'todos.updated'...")
	fmt.Println("Press Ctrl+C to stop")

	// Keep the connection alive
	select {}
}
