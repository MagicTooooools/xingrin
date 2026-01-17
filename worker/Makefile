.PHONY: build run test clean tidy

# Build the worker binary
build:
	go build -o bin/worker ./cmd/worker

# Run the worker
run:
	go run ./cmd/worker

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	go mod tidy
