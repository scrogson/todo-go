package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/scrogson/todo-go/internal/client"
	todov1 "github.com/scrogson/todo-go/pkg/todo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

// Version information - will be set by the build process
var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

func main() {
	// Check for version flag
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("Todo Client v%s (commit: %s, built: %s)\n", version, commit, buildTime)
		return
	}

	// Process command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Use recommended connection creation with NewClient
	conn, err := grpc.NewClient(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create client
	todoClient := client.NewTodoClient(todov1.NewTodoServiceClient(conn))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	command := os.Args[1]

	switch command {
	case "list":
		handleListTodos(ctx, todoClient)
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Title is required for add command")
			printUsage()
			return
		}
		title := strings.Join(os.Args[2:], " ")
		handleAddTodo(ctx, todoClient, title)
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: ID is required for delete command")
			printUsage()
			return
		}
		handleDeleteTodo(ctx, todoClient, os.Args[2])
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Error: ID and title are required for update command")
			printUsage()
			return
		}
		id := os.Args[2]
		title := strings.Join(os.Args[3:], " ")
		handleUpdateTodo(ctx, todoClient, id, title)
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Error: ID is required for complete command")
			printUsage()
			return
		}
		handleCompleteTodo(ctx, todoClient, os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

// CLI handler functions that use the client and format output

func handleListTodos(ctx context.Context, todoClient *client.TodoClient) {
	todos, err := todoClient.ListTodos(ctx)
	if err != nil {
		log.Fatalf("Could not list todos: %v", err)
	}

	if len(todos) == 0 {
		fmt.Println("No todos found.")
		return
	}

	fmt.Println("Todos:")
	for _, todo := range todos {
		status := " "
		if todo.Completed {
			status = "✓"
		}
		fmt.Printf("[%s] %s: %s\n", status, todo.Id, todo.Title)
	}
}

func handleAddTodo(ctx context.Context, todoClient *client.TodoClient, title string) {
	todo, err := todoClient.AddTodo(ctx, title)
	if err != nil {
		log.Fatalf("Could not add todo: %v", err)
	}
	fmt.Printf("Added todo: [%s] %s\n", todo.Id, todo.Title)
}

func handleDeleteTodo(ctx context.Context, todoClient *client.TodoClient, id string) {
	success, err := todoClient.DeleteTodo(ctx, id)
	if err != nil {
		log.Fatalf("Could not delete todo: %v", err)
	}
	if success {
		fmt.Println("Todo deleted successfully")
	} else {
		fmt.Println("Failed to delete todo")
	}
}

func handleUpdateTodo(ctx context.Context, todoClient *client.TodoClient, id, title string) {
	success, err := todoClient.UpdateTodo(ctx, id, title)
	if err != nil {
		log.Fatalf("Could not update todo: %v", err)
	}
	if success {
		fmt.Println("Todo updated successfully")
	} else {
		fmt.Println("Failed to update todo")
	}
}

func handleCompleteTodo(ctx context.Context, todoClient *client.TodoClient, id string) {
	success, err := todoClient.CompleteTodo(ctx, id)
	if err != nil {
		log.Fatalf("Could not complete todo: %v", err)
	}
	if success {
		fmt.Println("Todo marked as complete")
	} else {
		fmt.Println("Failed to mark todo as complete")
	}
}

// printUsage shows command line help
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  todo list                     - List all todos")
	fmt.Println("  todo add <title>              - Add a new todo")
	fmt.Println("  todo delete <id>              - Delete a todo by ID")
	fmt.Println("  todo update <id> <title>      - Update a todo's title")
	fmt.Println("  todo complete <id>            - Mark a todo as complete")
}
