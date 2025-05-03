# Makefile for building different binaries

.PHONY: all clean server client concurrent_client

all: server client concurrent_client

server:
	go build -o bin/server server.go 

client: 
	go build -o bin/client ./client/client.go

concurrent_client:
	go build -o bin/concurrent_client concurrent_client.go

run-server-pool: server
	./bin/server -pool

run-client: client
	./bin/client

run-concurrent-client: concurrent_client
	./bin/concurrent_client

clean:
	rm -rf bin/