.PHONY: help build install format lint test test-unit test-golden test-all test-coverage clean run tidy check ci all

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
		echo "golangci-lint not found. Skipping lint (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)"; \
	fi

test: test-all ## Run all tests (alias for test-all)

test-unit: ## Run unit tests (excludes golden tests)
	@echo "Running unit tests..."
	@go test -v ./internal/...

test-golden: build ## Run golden output tests
	@echo "Running golden tests..."
	@cd test && go test -v -run TestGolden

test-all: test-unit test-golden ## Run all tests (unit + golden)
	@echo "All tests complete"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -f test/data/*.duckdb test/data/*.duckdb.wal
	@rm -f test/golden/*/*.duckdb test/golden/*/*.duckdb.wal
	@go clean ./...
	@echo "Clean complete"

run: ## Run the application (example: make run ARGS="-e 'select 1'")
	@go run $(MAIN_PATH) $(ARGS)

tidy: ## Run go mod tidy
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies tidied"

check: format lint test-unit ## Run format, lint, and unit tests

ci: format lint test-all ## Full CI check (format, lint, all tests)
	@echo "CI checks complete"

all: clean format build test-all ## Clean, format, build, and test all
