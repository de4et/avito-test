setup:
	cp .env.example .env

docker-run:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	else \
		docker-compose up --build; \
	fi

docker-down:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	else \
		docker-compose down; \
	fi

all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

test:
	@go test ./... -v

clean:
	@rm -f main

.PHONY: all build run test clean docker-run docker-down
