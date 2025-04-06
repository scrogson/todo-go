package storage

import (
	"github.com/oklog/ulid/v2"
	todov1 "github.com/scrogson/todo-go/pkg/todo/v1"
)

// TodoStorage defines the interface for todo data storage
type TodoStorage interface {
	// Add create a new todo and returns its ID
	Add(title string) (*todov1.Todo, error)

	// Get returns a todo by ID
	Get(id ulid.ULID) (*todov1.Todo, bool)

	// List returns all todos sorted by ID
	// TODO: Add pagination
	List() ([]*todov1.Todo, error)

	// Update updates a todo's title
	Update(id ulid.ULID, title string) (bool, error)

	// Delete removes a todo by ID
	Delete(id ulid.ULID) (bool, error)

	// Complete marks a todo as completed
	Complete(id ulid.ULID) (bool, error)
}
