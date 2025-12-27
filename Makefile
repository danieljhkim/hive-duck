.PHONY: help build install format lint test clean run

# Binary name
BINARY_NAME=hive-duck
MAIN_PATH=./cmd/hive-duck

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: ./$(BINARY_NAME)"

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(MAIN_PATH)
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

format: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@go clean ./...
	@echo "Clean complete"

run: ## Run the application (example: make run ARGS="-e 'select 1'")
	@go run $(MAIN_PATH) $(ARGS)

tidy: ## Run go mod tidy
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies tidied"

check: format lint test ## Run format, lint, and test

all: clean format build test ## Clean, format, build, and test

