.PHONY: all build clean proto run-server run-client help deps test-deps test test-cover test-html bench bench-storage

# Binary names
SERVER_BINARY=server
CLIENT_BINARY=client

# Build directory
BUILD_DIR=bin

# Get the current directory
CURRENT_DIR=$(shell pwd)

# Version information
VERSION ?= $(shell git describe --tags --always --dirty || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD || echo "unknown")
BUILD_TIME ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

all: clean build

build: proto build-server build-client

# Install all dependencies
deps: proto-deps test-deps
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy

# Install protocol buffer dependencies
proto-deps:
	@echo "Checking Protocol Buffer compiler..."
	@if ! command -v protoc &> /dev/null; then \
		echo "protoc not found. Please install Protocol Buffers compiler:"; \
		echo "  macOS: brew install protobuf"; \
		echo "  Linux: apt-get install protobuf-compiler"; \
		exit 1; \
	fi
	@echo "Installing Protocol Buffer Go plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install test dependencies
test-deps:
	@echo "Installing test dependencies..."
	go get github.com/stretchr/testify/assert@latest

build-server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY) ./cmd/server

build-client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(CLIENT_BINARY) ./cmd/client

proto:
	@echo "Generating protocol buffer code..."
	protoc --proto_path=. --go_out=. --go-grpc_out=. proto/todo/v1/todo.proto

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Generate HTML coverage report
test-html:
	@echo "Generating HTML coverage report..."
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Run development versions
run-server:
	@echo "Running server..."
	go run ./cmd/server

run-client:
	@echo "Running client..."
	go run ./cmd/client $(ARGS)

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY)-linux-amd64 ./cmd/server
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(CLIENT_BINARY)-linux-amd64 ./cmd/client
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY)-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(CLIENT_BINARY)-darwin-amd64 ./cmd/client
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY)-darwin-arm64 ./cmd/server
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(CLIENT_BINARY)-darwin-arm64 ./cmd/client
	
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(SERVER_BINARY)-windows-amd64.exe ./cmd/server
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(CLIENT_BINARY)-windows-amd64.exe ./cmd/client

# Run benchmarks
bench:
	@echo "Running all benchmarks..."
	go test ./... -bench=. -benchmem

# Run storage benchmarks
bench-storage:
	@echo "Running storage benchmarks..."
	go test ./internal/storage -bench=. -benchmem

help:
	@echo "Available commands:"
	@echo "  make              - Build everything (same as 'make all')"
	@echo "  make build        - Build server and client binaries"
	@echo "  make build-server - Build only the server binary"
	@echo "  make build-client - Build only the client binary"
	@echo "  make proto        - Generate code from protobuf definitions"
	@echo "  make deps         - Install all dependencies"
	@echo "  make proto-deps   - Install Protocol Buffer dependencies"
	@echo "  make test-deps    - Install testing dependencies"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage summary"
	@echo "  make test-html    - Generate HTML coverage report"
	@echo "  make bench        - Run all benchmarks"
	@echo "  make bench-storage - Run storage benchmarks"
	@echo "  make run-server   - Run the server for development"
	@echo "  make run-client   - Run the client (use ARGS='list' for cli args)"
	@echo "  make build-all    - Build binaries for multiple platforms"
	@echo "  make help         - Show this help message" 