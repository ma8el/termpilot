# Termpilot Makefile

.PHONY: build test clean coverage lint all

# Default build target
all: test build

# Build the application
build:
	go build -o termpilot

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run individual package tests
test-db:
	go test -v ./db/...

test-models:
	go test -v ./models/...

test-cmd:
	go test -v ./cmd/...

test-ollama:
	go test -v ./ollamaclient/...

# Lint the code
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. To install, run:"; \
		echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	rm -f termpilot
	rm -f coverage.out
	rm -f coverage.html
	rm -f test.db

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install development tools
devtools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/stretchr/testify@latest 