package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"golang.org/x/sync/errgroup"
)

func main() {
	log.Println("Starting NATS microservice with errgroup worker pool")
	RunPoolServer()
}

const maxWorkers = 100 // Number of concurrent workers

// worker is the function that processes requests
func worker(req micro.Request) {
	// Process the request in the worker
	msg := string(req.Data())
	endpoint := req.Subject()
	log.Printf("Worker processing request from '%s': %s", endpoint, msg)

	// Simulate some work (e.g., processing)
	time.Sleep(100 * time.Millisecond)

	// Respond with the echo message including endpoint name
	resp := fmt.Sprintf("Response from %s: %s", endpoint, msg)
	if err := req.Respond([]byte(resp)); err != nil {
		log.Printf("Error responding to %s: %v", endpoint, err)
	} else {
		log.Printf("Worker responded to %s: %s", endpoint, resp)
	}
}

// RunPoolServer starts the NATS microservice with a worker pool
func RunPoolServer() {
	// Connect to NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Drain()

	// Start microservice
	svc, err := micro.AddService(nc, micro.Config{
		Name:    "echo_service", // Service name
		Version: "1.0.0",        // Service version
	})
	if err != nil {
		log.Fatalf("Error creating microservice: %v", err)
	}
	log.Println("Service started with worker pool: echo_service")

	// Create a context for the errgroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create an errgroup to manage worker goroutines
	g, ctx := errgroup.WithContext(ctx)

	// Create a buffered channel to manage incoming requests
	requestQueue := make(chan micro.Request, 100)

	// Start workers using errgroup
	for i := 0; i < maxWorkers; i++ {
		workerID := i // Capture loop variable
		g.Go(func() error {
			log.Printf("Worker #%d started", workerID)
			for {
				select {
				case req, ok := <-requestQueue:
					if !ok {
						// Channel closed, exit goroutine
						log.Printf("Worker #%d shutting down", workerID)
						return nil
					}

					log.Printf("Worker #%d processing request from '%s'", workerID, req.Subject())
					worker(req)
					log.Printf("Worker #%d completed request from '%s'", workerID, req.Subject())

				case <-ctx.Done():
					// Context canceled, exit goroutine
					log.Printf("Worker #%d shutting down due to context cancellation", workerID)
					return ctx.Err()
				}
			}
		})
	}

	// Create a handler function that sends requests to the queue
	createQueueHandler := func(endpointName string) micro.Handler {
		return micro.HandlerFunc(func(req micro.Request) {
			select {
			case requestQueue <- req:
				log.Printf("[%s] Request queued: %s", endpointName, string(req.Data()))
			case <-time.After(500 * time.Millisecond):
				log.Printf("[%s] Request queue full, rejecting: %s", endpointName, string(req.Data()))
				req.Respond([]byte("Error: Request queue is full, try again later"))
			}
		})
	}

	// Register endpoints
	endpoints := []string{"echo", "echo1"}
	for _, endpoint := range endpoints {
		err = svc.AddEndpoint(endpoint, createQueueHandler(endpoint))
		if err != nil {
			log.Fatalf("Error adding endpoint '%s': %v", endpoint, err)
		}
		log.Printf("Registered endpoint: %s", endpoint)
	}

	// Set up a goroutine to handle graceful shutdown when errors occur
	go func() {
		if err := g.Wait(); err != nil {
			log.Printf("Error in worker pool: %v", err)
			cancel() // Cancel context to signal all workers to shut down
		}
	}()

	// Block forever until context is canceled
	<-ctx.Done()
}
