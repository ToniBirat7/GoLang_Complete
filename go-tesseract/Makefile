# Go Tesseract OCR - Makefile

.PHONY: help build build-cli build-api clean test docker-build docker-up docker-down install-deps

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build both CLI and API binaries"
	@echo "  build-cli      - Build CLI binary"
	@echo "  build-api      - Build API binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  install-deps   - Download Go dependencies"

# Build both binaries
build: build-cli build-api

# Build CLI binary
build-cli:
	@echo "Building CLI..."
	@mkdir -p bin
	go build -o bin/ocr-cli ./cmd/ocr-cli

# Build API binary
build-api:
	@echo "Building API..."
	@mkdir -p bin
	go build -o bin/ocr-api ./cmd/ocr-api

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf uploads/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t go-tesseract-ocr:latest .

# Docker compose up
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# Run API locally
run-api:
	@echo "Starting API server..."
	go run ./cmd/ocr-api

# Run CLI with example
run-cli-example:
	@echo "Running CLI example..."
	go run ./cmd/ocr-cli -image example.png
