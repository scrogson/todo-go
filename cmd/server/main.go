package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	// Define flags
	storageType := flag.String("storage", "memory", "Storage type to use (memory or sqlite)")
	dbPath := flag.String("db", "todo.db", "Path to SQLite database file (only used with sqlite storage)")
	port := flag.Int("port", 50051, "Port to listen on")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create storage based on type
	var todoStorage storage.TodoStorage
	var sqliteStorage *storage.SQLiteStorage

	switch *storageType {
	case "memory":
		todoStorage = storage.NewInMemoryStorage()
		log.Printf("Using in-memory storage")
	case "sqlite":
		sqliteStorage, err = storage.NewSQLiteStorage(*dbPath)
		if err != nil {
			log.Fatalf("failed to create SQLite storage: %v", err)
		}
		todoStorage = sqliteStorage
		log.Printf("Using SQLite storage with database: %s", *dbPath)
	default:
		log.Fatalf("unknown storage type: %s", *storageType)
	}

	// Create server
	todoServer := server.NewTodoServer(todoStorage)

	// Create and start gRPC server
	grpcServer := grpc.NewServer()
	todov1.RegisterTodoServiceServer(grpcServer, todoServer)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to capture signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting Todo gRPC server v%s on :%d", version, *port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// Wait for termination signal
	select {
	case <-ctx.Done():
		// Context was canceled due to server error
	case sig := <-sigCh:
		log.Printf("Received signal: %v", sig)
	}

	// Graceful shutdown
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()

	// Close SQLite connection if used
	if sqliteStorage != nil {
		log.Println("Closing database connection...")
		if err := sqliteStorage.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	log.Println("Server shutdown complete")
}
