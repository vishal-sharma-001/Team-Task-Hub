.PHONY: help build run dev clean migrate fmt

help:
	@echo "Team Task Hub - Backend"
	@echo ""
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make dev      - Run without building (go run)"
	@echo "  make migrate  - Run database migrations"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make fmt      - Format code"

build:
	go build -o bin/team-task-hub ./cmd/team-task-hub/main.go

run: build
	./bin/team-task-hub

dev:
	go run ./cmd/team-task-hub/main.go

migrate:
	go run ./cmd/team-task-hub/main.go migrate

clean:
	rm -rf bin/

fmt:
	go fmt ./...

.DEFAULT_GOAL := help
