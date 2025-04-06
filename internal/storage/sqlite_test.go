package storage

import (
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStorage(t *testing.T) {
	// Use an in-memory database for testing
	dbPath := ":memory:"
	storage, err := NewSQLiteStorage(dbPath)
	require.NoError(t, err)
	defer storage.Close()

	// Test Add
	todo, err := storage.Add("Test Todo")
	require.NoError(t, err)
	assert.NotEmpty(t, todo.Id)
	assert.Equal(t, "Test Todo", todo.Title)
	assert.False(t, todo.Completed)

	// Test Get
	id, err := ulid.Parse(todo.Id)
	require.NoError(t, err)
	retrieved, exists := storage.Get(id)
	assert.True(t, exists)
	assert.Equal(t, todo.Id, retrieved.Id)
	assert.Equal(t, todo.Title, retrieved.Title)
	assert.Equal(t, todo.Completed, retrieved.Completed)

	// Test non-existent Get
	invalidID := ulid.MustNew(1, nil)
	_, exists = storage.Get(invalidID)
	assert.False(t, exists)

	// Test List
	todos, err := storage.List()
	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, todo.Id, todos[0].Id)

	// Test Update
	updated, err := storage.Update(id, "Updated Todo")
	require.NoError(t, err)
	assert.True(t, updated)

	retrieved, exists = storage.Get(id)
	assert.True(t, exists)
	assert.Equal(t, "Updated Todo", retrieved.Title)

	// Test update with empty title
	updated, err = storage.Update(id, "")
	assert.Error(t, err)
	assert.False(t, updated)

	// Test Complete
	completed, err := storage.Complete(id)
	require.NoError(t, err)
	assert.True(t, completed)

	retrieved, exists = storage.Get(id)
	assert.True(t, exists)
	assert.True(t, retrieved.Completed)

	// Test Delete
	deleted, err := storage.Delete(id)
	require.NoError(t, err)
	assert.True(t, deleted)

	_, exists = storage.Get(id)
	assert.False(t, exists)

	// Test Delete non-existent
	deleted, err = storage.Delete(id)
	require.NoError(t, err)
	assert.False(t, deleted)
}

func TestSQLiteStorageFile(t *testing.T) {
	// Use a temporary file for testing
	tempFile, err := os.CreateTemp("", "todo-sqlite-test-*.db")
	require.NoError(t, err)
	tempFile.Close()
	dbPath := tempFile.Name()
	defer os.Remove(dbPath)

	// Create and test with a file-based SQLite database
	storage, err := NewSQLiteStorage(dbPath)
	require.NoError(t, err)
	defer storage.Close()

	// Add a todo
	todo, err := storage.Add("File-based Todo")
	require.NoError(t, err)
	assert.NotEmpty(t, todo.Id)

	// Close the connection
	err = storage.Close()
	require.NoError(t, err)

	// Reopen the database and verify persistence
	reopenedStorage, err := NewSQLiteStorage(dbPath)
	require.NoError(t, err)
	defer reopenedStorage.Close()

	id, err := ulid.Parse(todo.Id)
	require.NoError(t, err)

	retrieved, exists := reopenedStorage.Get(id)
	assert.True(t, exists)
	assert.Equal(t, todo.Id, retrieved.Id)
	assert.Equal(t, todo.Title, retrieved.Title)
}
