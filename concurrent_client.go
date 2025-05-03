package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS server at", nats.DefaultURL)

	// Number of concurrent requests to send
	const numRequests = 1000

	// Available endpoints
	endpoints := []string{"echo", "echo1"}

	// Create a wait group to wait for all requests to complete
	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Send multiple requests concurrently
	log.Printf("Sending %d concurrent requests...", numRequests)
	for i := 1; i <= numRequests; i++ {
		go func(requestNum int) {
			defer wg.Done()

			// Choose endpoint - alternate between echo and echo1
			endpointIndex := (requestNum - 1) % len(endpoints)
			endpoint := endpoints[endpointIndex]

			// Prepare payload
			payload := fmt.Sprintf("Request %d", requestNum)

			// Log the request
			log.Printf("[Client] Sending request #%d to %s: %s",
				requestNum, endpoint, payload)

			// Send request
			startTime := time.Now()
			resp, err := nc.Request(endpoint, []byte(payload), 15*time.Second)

			if err != nil {
				log.Printf("[Client] Error on request #%d to %s: %v",
					requestNum, endpoint, err)
				return
			}

			elapsed := time.Since(startTime)
			log.Printf("[Client] Received response from %s, request #%d after %v: %s",
				endpoint, requestNum, elapsed, string(resp.Data))
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	log.Println("All requests completed")
}
