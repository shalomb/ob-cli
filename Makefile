# ob-cli Makefile
# Fast Obsidian/Tips CLI with frecency and async git sync

.PHONY: build test test-integration clean install deps lint format

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_USER ?= $(shell whoami)
BUILD_HOST ?= $(shell hostname)

# Build flags
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildUser=$(BUILD_USER) -X main.BuildHost=$(BUILD_HOST)"

# Build the CLI tool
build:
	@echo "Building ob-cli..."
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	go build $(LDFLAGS) -o bin/ob-cli ./cmd/ob-cli
	@echo "Build complete: bin/ob-cli"

# Run tests
test:
	go test -v ./...

# Run integration tests (with real file operations)
test-integration:
	go test -v -tags=integration ./test/integration/...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Install to ~/.local/bin
install: build
	cp bin/ob-cli ~/.local/bin/ob-cli
	chmod +x ~/.local/bin/ob-cli

# Install dependencies
deps:
	go mod tidy
	go mod download

# Lint code
lint:
	golangci-lint run

# Format code
format:
	go fmt ./...
	goimports -w .

# Development mode - build and run
dev: build
	./bin/ob-cli

# Run with specific mode
dev-tips: build
	./bin/ob-cli --mode=tips

dev-obsidian: build
	./bin/ob-cli --mode=obsidian

# Generate mocks
mocks:
	go generate ./...

# Run BDD tests
bdd:
	ginkgo run ./test/bdd/...

# Coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"