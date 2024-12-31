package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type UserEvent struct {
	ID      string          `json:id`
	Type    string          `json:type`
	Repo    Repo            `json:repo`
	Payload json.RawMessage `json:payload`
}

type Repo struct {
	Name string `json: name`
}

type PushEvent struct {
	Size int
}

func main() {
	if len(os.Args) < 2 {
		os.Exit(0)
	}
	username := os.Args[1]
	fetchActivity(username)
}

func fetchActivity(username string) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	var userEvents []UserEvent
	err = json.Unmarshal(body, &userEvents)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range userEvents {
		switch event.Type {
		case "PushEvent":
			var pushEvent PushEvent
			err := json.Unmarshal(event.Payload, &pushEvent)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("- Pushed %d commits to %s\n", pushEvent.Size, event.Repo.Name)
		default:
			fmt.Printf("- Missing handling for event type \"%s\"\n", event.Type)
		}
	}
}
