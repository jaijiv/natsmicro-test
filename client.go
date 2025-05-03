package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type GreetRequest struct {
	Name string `json:"name"`
}

type GreetResponse struct {
	Greeting string `json:"greeting"`
}

func main() {
	// Connect to NATS with options
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS server at", nats.DefaultURL)

	// Send a request to the service endpoint using direct NATS request
	subject := "echo"
	log.Printf("Sending request to %s...", subject)
	respMsg, err := nc.Request(subject, []byte("hello"), 20*time.Second)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	fmt.Printf("\nReceived response: %s\n", string(respMsg.Data))
}
