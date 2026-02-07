.PHONY: build run test check-ddd lint clean check-dto-selftest

# Build the server binary
build:
	go build -o bin/server ./cmd/server

# Run the server
run:
	go run ./cmd/server

# DDD + DTO boundary guard
check-ddd:
	bash ./scripts/check-ddd-boundaries.sh
	bash ./scripts/check-dto-boundaries.sh

# Local optional guard self-test
check-dto-selftest:
	bash ./scripts/check-dto-boundaries-selftest.sh


# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Generate mocks (if needed)
generate:
	go generate ./...
