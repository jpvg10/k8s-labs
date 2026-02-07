package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func writeToFile(content string) {
	dir := filepath.Join("/", "usr", "src", "app", "files")
	filepath := filepath.Join(dir, "output.txt")

	os.MkdirAll(dir, 0755)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
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
		output := fmt.Sprintf("[TODO CREATED] ID: %d, Description: %s, Done: %v\n", todo.Id, todo.Description, todo.Done)
		fmt.Print(output)
		writeToFile(output)
	})

	nc.Subscribe("todos.updated", func(msg *nats.Msg) {
		var todo Todo
		if err := json.Unmarshal(msg.Data, &todo); err != nil {
			fmt.Printf("Error unmarshaling updated todo: %v\n", err)
			return
		}
		output := fmt.Sprintf("[TODO UPDATED] ID: %d, Description: %s, Done: %v\n", todo.Id, todo.Description, todo.Done)
		fmt.Print(output)
		writeToFile(output)
	})

	fmt.Println("Listening for NATS messages on 'todos.created' and 'todos.updated'...")
	fmt.Println("Press Ctrl+C to stop")

	// Keep the connection alive
	select {}
}
