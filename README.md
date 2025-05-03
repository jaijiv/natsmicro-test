# NATS Microservice Example with Worker Pool

This repository demonstrates an efficient worker pool implementation for NATS microservices in Go.

## Implementation with errgroup

The server implements a NATS microservice with an errgroup-based worker pool. This implementation offers significant advantages:

1. **Better error handling**: Uses errgroup for proper error propagation between workers
2. **Graceful shutdown**: Context cancellation allows clean shutdown of all workers
3. **Centralized queue**: Single request queue for better load distribution
4. **Scalable design**: Easily adjust number of workers (currently set to 100) based on load requirements
5. **Reliable processing**: Prevents race conditions and request timeouts

## How to Run

Build and run the server:
```
go run server.go
```

Run a client:
```
go run concurrent_client.go
```

Use the NATS CLI to test individual endpoints:
```
nats req echo "Hello World"
nats req echo1 "Hello World"
```

## Implementation Details

The errgroup-based worker pool implementation:

1. Creates a central request queue for all incoming requests
2. Manages worker goroutines with errgroup for better error handling
3. Uses context cancellation for graceful shutdown
4. Supports multiple endpoints ("echo" and "echo1") sharing the same worker pool
5. Handles high concurrency with 100 workers and request queue buffer of 100
6. Provides request timeout handling to prevent blocking

This approach follows idiomatic Go patterns for concurrent programming with proper error handling and resource management.