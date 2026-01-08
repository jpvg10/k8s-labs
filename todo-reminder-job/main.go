package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	todoBackendUrl := os.Getenv("TODO_BACKEND_URL")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "https://en.wikipedia.org/wiki/Special:Random", nil)
	req.Header.Set("User-Agent", "GoBot/1.0")
	res, err := client.Do(req)

	if err != nil {
		fmt.Println("Error getting Wikipedia article")
	} else {
		url, _ := res.Location()
		desc := fmt.Sprintf("Read: %s", url)
		jsonValue, _ := json.Marshal(map[string]any{"description": desc})

		res, err := http.Post(fmt.Sprintf("%s/todos", todoBackendUrl), "application/json", bytes.NewBuffer(jsonValue))
		if err != nil || res.StatusCode >= 400 {
			fmt.Println("Error posting the Todo")
			return
		}

		fmt.Printf("Created todo:\n\t%s\n", desc)
	}
}
