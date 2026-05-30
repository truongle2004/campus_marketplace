# Simple Makefile for a Go project
COVERAGE_DIR := coverage


# Detect OS
ifeq ($(OS),Windows_NT)
	OPEN_CMD := start
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		OPEN_CMD := open
	else
		OPEN_CMD := xdg-open
	endif
endif

# Seed dev data
seed:
	@go run cmd/seed/main.go

# Generate Wire DI code
wire:
	@go run github.com/google/wire/cmd/wire@latest ./internal/di

# Run GORM auto-migrations
migrate:
	@go run cmd/migrate/main.go

# Build the application
all: build test

build:
	@echo "Building..."


	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

docker-build:
	echo "Building Docker image..."
	docker compose build;

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	@go test ./... -coverprofile=$(COVERAGE_DIR)/coverage.out
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"
	@echo "Opening coverage report..."
	@$(OPEN_CMD) $(COVERAGE_DIR)/coverage.html || echo "Please open it manually."

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest migrate wire seed
