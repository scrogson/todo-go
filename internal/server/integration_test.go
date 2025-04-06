package server

import (
	"context"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/scrogson/todo-go/internal/storage"
	todov1 "github.com/scrogson/todo-go/pkg/todo/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTodoServerIntegration tests the TodoServer with a real storage implementation
func TestTodoServerIntegration(t *testing.T) {
	// Create a real storage implementation
	todoStorage := storage.NewInMemoryStorage()

	// Create the server with real storage
	server := NewTodoServer(todoStorage)

	ctx := context.Background()

	// Test adding a todo
	addResp, err := server.AddTodo(ctx, &todov1.AddTodoRequest{
		Title: "Integration Test Todo",
	})
	require.NoError(t, err)
	require.NotNil(t, addResp)
	require.NotNil(t, addResp.Todo)
	require.Equal(t, "Integration Test Todo", addResp.Todo.Title)
	require.False(t, addResp.Todo.Completed)
	require.NotEmpty(t, addResp.Todo.Id)

	todoID := addResp.Todo.Id

	// Test listing todos
	listResp, err := server.ListTodos(ctx, &todov1.ListTodosRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Len(t, listResp.Todos, 1)
	require.Equal(t, todoID, listResp.Todos[0].Id)

	// Test updating a todo
	parsedID, err := ulid.Parse(todoID)
	require.NoError(t, err)

	updateResp, err := server.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
		Id:    todoID,
		Title: "Updated Integration Test Todo",
	})
	require.NoError(t, err)
	require.NotNil(t, updateResp)
	require.True(t, updateResp.Success)

	// Verify the update via storage directly
	updatedTodo, exists := todoStorage.Get(parsedID)
	require.True(t, exists)
	require.Equal(t, "Updated Integration Test Todo", updatedTodo.Title)

	// Test completing a todo
	completeResp, err := server.CompleteTodo(ctx, &todov1.CompleteTodoRequest{
		Id: todoID,
	})
	require.NoError(t, err)
	require.NotNil(t, completeResp)
	require.True(t, completeResp.Success)

	// Verify completion via storage directly
	completedTodo, exists := todoStorage.Get(parsedID)
	require.True(t, exists)
	require.True(t, completedTodo.Completed)

	// Test deleting a todo
	deleteResp, err := server.DeleteTodo(ctx, &todov1.DeleteTodoRequest{
		Id: todoID,
	})
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
	require.True(t, deleteResp.Success)

	// Verify deletion via storage directly
	_, exists = todoStorage.Get(parsedID)
	require.False(t, exists)

	// Test listing after deletion
	listResp, err = server.ListTodos(ctx, &todov1.ListTodosRequest{})
	require.NoError(t, err)
	require.NotNil(t, listResp)
	require.Empty(t, listResp.Todos)
}

// TestTodoServerIntegrationErrors tests error scenarios in the TodoServer with a real storage
func TestTodoServerIntegrationErrors(t *testing.T) {
	// Create a real storage implementation
	todoStorage := storage.NewInMemoryStorage()

	// Create the server with real storage
	server := NewTodoServer(todoStorage)

	ctx := context.Background()

	// Test adding with empty title
	_, err := server.AddTodo(ctx, &todov1.AddTodoRequest{
		Title: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title cannot be empty")

	// Test operations with invalid ID
	invalidID := "invalid-id"

	// Update with invalid ID
	_, err = server.UpdateTodo(ctx, &todov1.UpdateTodoRequest{
		Id:    invalidID,
		Title: "Updated Title",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ID")

	// Delete with invalid ID
	_, err = server.DeleteTodo(ctx, &todov1.DeleteTodoRequest{
		Id: invalidID,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ID")

	// Complete with invalid ID
	_, err = server.CompleteTodo(ctx, &todov1.CompleteTodoRequest{
		Id: invalidID,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ID")
}
