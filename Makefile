.PHONY: all build clean proto run-server run-client help

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

help:
	@echo "Available commands:"
	@echo "  make              - Build everything (same as 'make all')"
	@echo "  make build        - Build server and client binaries"
	@echo "  make build-server - Build only the server binary"
	@echo "  make build-client - Build only the client binary"
	@echo "  make proto        - Generate code from protobuf definitions"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make run-server   - Run the server for development"
	@echo "  make run-client   - Run the client (use ARGS='list' for cli args)"
	@echo "  make build-all    - Build binaries for multiple platforms"
	@echo "  make help         - Show this help message" 