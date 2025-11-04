.PHONY: help build run test test-coverage lint clean docker-build docker-run docker-stop install-tools

# Variables
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)
DOCKER_IMAGE=go-app
DOCKER_TAG=latest
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Default target
.DEFAULT_GOAL := help

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags="-s -w" -o $(BINARY_PATH) ./cmd/server
	@echo "Build complete: $(BINARY_PATH)"

build-linux: ## Build for Linux
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_PATH)-linux-amd64 ./cmd/server

build-darwin: ## Build for macOS
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_PATH)-darwin-amd64 ./cmd/server

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY_PATH)-windows-amd64.exe ./cmd/server

build-all: build-linux build-darwin build-windows ## Build for all platforms

run: build ## Run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

dev: ## Run with hot reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

test: ## Run tests
	@echo "Running tests..."
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-short: ## Run short tests
	@echo "Running short tests..."
	go test -v -short ./...

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

lint: ## Run linters
	@echo "Running linters..."
	go vet ./...
	golangci-lint run --timeout=5m

security: ## Run security scanner
	@echo "Running security scanner..."
	gosec -fmt=json -out=security-report.json ./...
	gosec ./...

format: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w $(GO_FILES)

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f security-report.json
	go clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: docker-build ## Run Docker container
	@echo "Running Docker container..."
	docker run -d \
		--name $(DOCKER_IMAGE) \
		-p 8080:8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_IMAGE) || true
	docker rm $(DOCKER_IMAGE) || true

docker-compose-up: ## Start services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services with docker-compose..."
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

ci: deps lint test ## Run CI pipeline locally

all: clean deps lint test build ## Run all tasks

.PHONY: all
