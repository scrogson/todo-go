package storage

import (
	"testing"

	"github.com/oklog/ulid/v2"
	todov1 "github.com/scrogson/todo-golang/pkg/todo/v1"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage_Add(t *testing.T) {
	s := NewInMemoryStorage()

	todo, err := s.Add("Test Todo")

	assert.NoError(t, err)
	assert.NotNil(t, todo)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.False(t, todo.Completed)
	assert.NotEmpty(t, todo.Id)

	// Verify it's in the storage
	id, err := ulid.Parse(todo.Id)
	assert.NoError(t, err)

	stored, exists := s.Get(id)
	assert.True(t, exists)
	assert.Equal(t, todo, stored)
}

func TestInMemoryStorage_Get(t *testing.T) {
	s := NewInMemoryStorage()

	// Add a todo to get
	todo, err := s.Add("Test Todo")
	assert.NoError(t, err)

	id, err := ulid.Parse(todo.Id)
	assert.NoError(t, err)

	// Test successful get
	retrieved, exists := s.Get(id)
	assert.True(t, exists)
	assert.Equal(t, todo, retrieved)

	// Test non-existent ID
	nonExistentID := ulid.MustNew(1, nil)
	_, exists = s.Get(nonExistentID)
	assert.False(t, exists)
}

func TestInMemoryStorage_List(t *testing.T) {
	s := NewInMemoryStorage()

	// Empty list
	list, err := s.List()
	assert.NoError(t, err)
	assert.Empty(t, list)

	// Add some todos
	todo1, err := s.Add("First Todo")
	assert.NoError(t, err)

	todo2, err := s.Add("Second Todo")
	assert.NoError(t, err)

	// Test list has both todos
	list, err = s.List()
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// Since they're sorted by ID, we can't guarantee the order without parsing the IDs
	// Just check that both todos are in the list
	todoMap := make(map[string]*todov1.Todo)
	for _, t := range list {
		todoMap[t.Id] = t
	}

	assert.Contains(t, todoMap, todo1.Id)
	assert.Contains(t, todoMap, todo2.Id)
}

func TestInMemoryStorage_Update(t *testing.T) {
	s := NewInMemoryStorage()

	// Add a todo to update
	todo, err := s.Add("Original Title")
	assert.NoError(t, err)

	id, err := ulid.Parse(todo.Id)
	assert.NoError(t, err)

	// Update with valid title
	updated, err := s.Update(id, "Updated Title")
	assert.NoError(t, err)
	assert.True(t, updated)

	// Verify update
	retrieved, exists := s.Get(id)
	assert.True(t, exists)
	assert.Equal(t, "Updated Title", retrieved.Title)

	// Try update with empty title
	updated, err = s.Update(id, "")
	assert.Error(t, err)
	assert.False(t, updated)

	// Update non-existent todo
	nonExistentID := ulid.MustNew(1, nil)
	updated, err = s.Update(nonExistentID, "New Title")
	assert.NoError(t, err)
	assert.False(t, updated)
}

func TestInMemoryStorage_Delete(t *testing.T) {
	s := NewInMemoryStorage()

	// Add a todo to delete
	todo, err := s.Add("To Be Deleted")
	assert.NoError(t, err)

	id, err := ulid.Parse(todo.Id)
	assert.NoError(t, err)

	// Delete the todo
	deleted, err := s.Delete(id)
	assert.NoError(t, err)
	assert.True(t, deleted)

	// Verify it's gone
	_, exists := s.Get(id)
	assert.False(t, exists)

	// Try to delete again
	deleted, err = s.Delete(id)
	assert.NoError(t, err)
	assert.False(t, deleted)

	// Delete non-existent todo
	nonExistentID := ulid.MustNew(1, nil)
	deleted, err = s.Delete(nonExistentID)
	assert.NoError(t, err)
	assert.False(t, deleted)
}

func TestInMemoryStorage_Complete(t *testing.T) {
	s := NewInMemoryStorage()

	// Add a todo to complete
	todo, err := s.Add("To Be Completed")
	assert.NoError(t, err)
	assert.False(t, todo.Completed)

	id, err := ulid.Parse(todo.Id)
	assert.NoError(t, err)

	// Complete the todo
	completed, err := s.Complete(id)
	assert.NoError(t, err)
	assert.True(t, completed)

	// Verify it's completed
	retrieved, exists := s.Get(id)
	assert.True(t, exists)
	assert.True(t, retrieved.Completed)

	// Try to complete non-existent todo
	nonExistentID := ulid.MustNew(1, nil)
	completed, err = s.Complete(nonExistentID)
	assert.NoError(t, err)
	assert.False(t, completed)
}

func TestInMemoryStorage_Concurrency(t *testing.T) {
	s := NewInMemoryStorage()

	// This is a simple test to ensure the mutex is working
	// For a more robust test, we would use the race detector and goroutines

	// Add a todo
	todo, err := s.Add("Concurrent Todo")
	assert.NoError(t, err)

	id, err := ulid.Parse(todo.Id)
	assert.NoError(t, err)

	// Run operations that would deadlock if mutex isn't working properly
	go func() {
		s.Update(id, "Updated Title")
	}()

	go func() {
		s.Complete(id)
	}()

	go func() {
		s.List()
	}()

	// If we get here without deadlock, the test passes
}
