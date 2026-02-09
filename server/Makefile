.PHONY: build run test check-ddd check-architecture check-layer-dependencies check-domain-purity check-no-legacy-layer check-dto-selftest check-repo-selftest check-common-httpdto sync-common-httpdto lint clean deps fmt generate

# Build the server binary
build:
	go build -o bin/server ./cmd/server

# Run the server
run:
	go run ./cmd/server

# Architecture guard (DDD-Strict + DTO + Repository + HTTPDTO)
check-architecture:
	bash ./scripts/check-no-legacy-layer.sh
	bash ./scripts/check-layer-dependencies.sh
	bash ./scripts/check-domain-purity.sh
	bash ./scripts/check-ddd-boundaries.sh
	bash ./scripts/check-dto-boundaries.sh
	bash ./scripts/check-repository-boundaries.sh
	bash ./scripts/check-common-httpdto.sh

check-layer-dependencies:
	bash ./scripts/check-layer-dependencies.sh

check-domain-purity:
	bash ./scripts/check-domain-purity.sh

check-no-legacy-layer:
	bash ./scripts/check-no-legacy-layer.sh

# Keep backward-compatible entrypoint
check-ddd:
	$(MAKE) check-architecture

# Local optional DTO guard self-test
check-dto-selftest:
	bash ./scripts/check-dto-boundaries-selftest.sh

# Local optional repository guard self-test
check-repo-selftest:
	bash ./scripts/check-repository-boundaries-selftest.sh

# Check common_http template consistency
check-common-httpdto:
	bash ./scripts/check-common-httpdto.sh

# Sync common_http files from template
sync-common-httpdto:
	bash ./scripts/sync-common-httpdto.sh

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
