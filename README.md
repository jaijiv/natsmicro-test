# NATS Microservice Example with Worker Pool

This repository demonstrates different worker pool implementations for NATS microservices in Go.

## Pool Server Implementation

The `poolserver.go` file implements a NATS microservice with a `sync.Pool` of worker goroutines. This implementation offers advantages over a fixed worker pool:

1. **Better resource utilization**: Workers are returned to the pool after handling a request
2. **Automatic scaling**: The pool can grow or shrink based on demand
3. **Improved performance**: Reusing goroutines reduces allocation overhead
4. **Reduced contention**: Workers are independently managed and don't block each other

## How to Run

Build all binaries:
```
make all
```

Run the server with sync.Pool implementation:
```
make run-server-pool
```

Run a client:
```
make run-client
```

Run a concurrent client (sends multiple requests):
```
make run-concurrent-client
```

## Implementation Details

The sync.Pool implementation in `poolserver.go`:

1. Creates a pool of worker channels
2. Pre-allocates a fixed number of worker goroutines
3. Returns worker channels to the pool after handling a request
4. Automatically handles request distribution

This approach maintains the benefits of a worker pool while allowing more efficient resource utilization.