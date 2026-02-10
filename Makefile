.PHONY: run build db-up db-down

## run: Run the application
run:
	go run ./cmd/server

## build: Build the application
build:
	go build -o bin/server ./cmd/server

## db-up: Start PostgreSQL with Docker
db-up:
	docker compose up -d postgres

## db-down: Stop PostgreSQL
db-down:
	docker compose down

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'