run:
	go run ./cmd/server

build:
	go build -o bin/asky ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-reup:
	go run ./cmd/migrate down
	go run ./cmd/migrate up

migrate-version:
	go run ./cmd/migrate version