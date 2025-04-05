# Todo CLI Application

[![Go CI](https://github.com/scrogson/todo-golang/actions/workflows/ci.yml/badge.svg)](https://github.com/scrogson/todo-golang/actions/workflows/ci.yml)

A simple Todo application built with Go, Protocol Buffers, and gRPC.

## Features

- Create, read, update, and delete todo items
- Mark todos as completed
- Command-line interface
- gRPC server for the backend
- Uses ULIDs for identifiers (time-ordered and sortable)

## Requirements

- Go 1.20+
- Protocol Buffers compiler (`protoc`)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

## Installation

### Installing Dependencies

```bash
make deps
```

### Building the Application

Clone the repository:

```bash
git clone https://github.com/scrogson/todo-golang.git
cd todo-golang
```

Build the server and client:

```bash
make build
```

This will create the binaries in the `bin/` directory.

## Usage

### Starting the Server

```bash
make run-server
# Or run the binary directly
./bin/server
```

### Using the Client

```bash
# Add a new todo
make run-client ARGS='add "Buy groceries"'
# Or run the binary directly
./bin/client add "Buy groceries"

# List all todos
make run-client ARGS='list'
# Or
./bin/client list

# Complete a todo (replace ID with actual ULID)
make run-client ARGS='complete 01FZGTA3JVT7RX870HAGBDXX9N'
# Or
./bin/client complete 01FZGTA3JVT7RX870HAGBDXX9N

# Update a todo (replace ID with actual ULID)
make run-client ARGS='update 01FZGTA3JVT7RX870HAGBDXX9N "Buy organic groceries"'
# Or
./bin/client update 01FZGTA3JVT7RX870HAGBDXX9N "Buy organic groceries"

# Delete a todo (replace ID with actual ULID)
make run-client ARGS='delete 01FZGTA3JVT7RX870HAGBDXX9N'
# Or
./bin/client delete 01FZGTA3JVT7RX870HAGBDXX9N
```

## Project Structure

```
todo-golang/
├── cmd/                # Command-line applications
│   ├── client/         # CLI client
│   └── server/         # gRPC server
├── internal/           # Private application code
│   ├── client/         # Client library
│   ├── server/         # Server implementation
│   └── storage/        # Data storage interface and implementations
├── pkg/                # Public libraries
│   └── todo/
│       └── v1/         # Generated Protocol Buffer code
└── proto/              # Protocol Buffer definitions
    └── todo/
        └── v1/         # Todo service definition
```

## Development

### Regenerating Protocol Buffer Code

If you make changes to the Protocol Buffer definitions, regenerate the Go code:

```bash
make proto
```

### Building for Multiple Platforms

```bash
make build-all
```

This creates binaries for Linux, macOS, and Windows in the `bin/` directory.

### Running Tests

```bash
make test
```

### Code Quality

#### Running Linters

```bash
# Install linters
make lint-install

# Run linters
make lint

# Auto-fix linting issues
make lint-fix
```

### Benchmarking

```bash
# Run all benchmarks
make bench

# Run only storage benchmarks
make bench-storage
```

## Continuous Integration

This project uses GitHub Actions for continuous integration. The CI pipeline includes:

- **Linting**: Ensures code follows best practices and style guidelines
- **Testing**: Runs the test suite with coverage reporting
- **Building**: Compiles the application for verification
- **Benchmarking**: Runs performance benchmarks
- **Releasing**: Automatically creates releases with binaries for multiple platforms when tags are pushed

To trigger a new release, create and push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## License

MIT
