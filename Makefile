run:
	go run ./cmd/server

build:
	go build -o bin/asky ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...
