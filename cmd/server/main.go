package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/scrogson/todo-golang/internal/server"
	"github.com/scrogson/todo-golang/internal/storage"
	todov1 "github.com/scrogson/todo-golang/pkg/todo/v1"
	"google.golang.org/grpc"
)

// Version information - will be set by the build process
var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

func main() {
	// Print version info if requested
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Todo Server v%s (commit: %s, built: %s)\n", version, commit, buildTime)
		return
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create storage
	todoStorage := storage.NewInMemoryStorage()

	// Create server
	todoServer := server.NewTodoServer(todoStorage)

	// Create and start gRPC server
	grpcServer := grpc.NewServer()
	todov1.RegisterTodoServiceServer(grpcServer, todoServer)

	log.Printf("Starting Todo gRPC server v%s on :50051", version)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
