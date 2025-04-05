package hello

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"github.com/scrogson/todo-golang/internal/client"
	"github.com/scrogson/todo-golang/internal/server"
	"github.com/scrogson/todo-golang/internal/storage"
	todov1 "github.com/scrogson/todo-golang/pkg/todo/v1"
)

// setupGRPCServer creates an in-memory gRPC server using bufconn for testing
func setupGRPCServer(t *testing.T) (*grpc.ClientConn, func()) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	// Create a gRPC server
	s := grpc.NewServer()

	// Create storage and server
	todoStorage := storage.NewInMemoryStorage()
	todoServer := server.NewTodoServer(todoStorage)

	// Register the server
	todov1.RegisterTodoServiceServer(s, todoServer)

	// Start the server
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Failed to serve: %v", err)
		}
	}()

	// Create a client connection
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return conn, func() {
		conn.Close()
		s.Stop()
	}
}

// TestEndToEndIntegration tests the flow from client to server to storage and back
func TestEndToEndIntegration(t *testing.T) {
	// Set up the server and client connection
	conn, cleanup := setupGRPCServer(t)
	defer cleanup()

	// Create the gRPC client
	grpcClient := todov1.NewTodoServiceClient(conn)

	// Create the Todo client
	todoClient := client.NewTodoClient(grpcClient)

	ctx := context.Background()

	// 1. Test adding a todo
	todo, err := todoClient.AddTodo(ctx, "End-to-End Test Todo")
	require.NoError(t, err)
	require.NotNil(t, todo)
	require.Equal(t, "End-to-End Test Todo", todo.Title)
	require.False(t, todo.Completed)
	require.NotEmpty(t, todo.Id)

	todoID := todo.Id

	// 2. Test listing todos
	todos, err := todoClient.ListTodos(ctx)
	require.NoError(t, err)
	require.Len(t, todos, 1)
	require.Equal(t, todoID, todos[0].Id)

	// 3. Test updating a todo
	updated, err := todoClient.UpdateTodo(ctx, todoID, "Updated E2E Test Todo")
	require.NoError(t, err)
	require.True(t, updated)

	// 4. Verify the update
	todos, err = todoClient.ListTodos(ctx)
	require.NoError(t, err)
	require.Len(t, todos, 1)
	require.Equal(t, "Updated E2E Test Todo", todos[0].Title)

	// 5. Test completing a todo
	completed, err := todoClient.CompleteTodo(ctx, todoID)
	require.NoError(t, err)
	require.True(t, completed)

	// 6. Verify the completion
	todos, err = todoClient.ListTodos(ctx)
	require.NoError(t, err)
	require.Len(t, todos, 1)
	require.True(t, todos[0].Completed)

	// 7. Test deleting a todo
	deleted, err := todoClient.DeleteTodo(ctx, todoID)
	require.NoError(t, err)
	require.True(t, deleted)

	// 8. Verify the deletion
	todos, err = todoClient.ListTodos(ctx)
	require.NoError(t, err)
	require.Empty(t, todos)
}

// TestEndToEndErrors tests error scenarios in the client-server integration
func TestEndToEndErrors(t *testing.T) {
	// Set up the server and client connection
	conn, cleanup := setupGRPCServer(t)
	defer cleanup()

	// Create the gRPC client
	grpcClient := todov1.NewTodoServiceClient(conn)

	// Create the Todo client
	todoClient := client.NewTodoClient(grpcClient)

	ctx := context.Background()

	// 1. Test adding with empty title
	_, err := todoClient.AddTodo(ctx, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot be empty")

	// 2. Test operations with non-existent ID
	nonExistentID := "01J3VC7K7C9P9M2H6T5QDNBGBZ"

	// Update with non-existent ID
	updated, err := todoClient.UpdateTodo(ctx, nonExistentID, "Updated Title")
	assert.NoError(t, err) // The error is translated to a boolean result
	assert.False(t, updated)

	// Delete with non-existent ID
	deleted, err := todoClient.DeleteTodo(ctx, nonExistentID)
	assert.NoError(t, err) // The error is translated to a boolean result
	assert.False(t, deleted)

	// Complete with non-existent ID
	completed, err := todoClient.CompleteTodo(ctx, nonExistentID)
	assert.NoError(t, err) // The error is translated to a boolean result
	assert.False(t, completed)

	// 3. Test operations with invalid ID
	invalidID := "invalid-id"

	// Update with invalid ID
	_, err = todoClient.UpdateTodo(ctx, invalidID, "Updated Title")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ID")

	// Delete with invalid ID
	_, err = todoClient.DeleteTodo(ctx, invalidID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ID")

	// Complete with invalid ID
	_, err = todoClient.CompleteTodo(ctx, invalidID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ID")
}
