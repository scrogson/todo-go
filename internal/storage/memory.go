package storage

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	todov1 "github.com/scrogson/todo-go/pkg/todo/v1"
)

// InMemoryStorage implements TodoStorage interface with in-memory storage
type InMemoryStorage struct {
	mu    sync.RWMutex
	todos map[ulid.ULID]*todov1.Todo
	rnd   *rand.Rand
}

// NewInMemoryStorage creates a new in-memory storage instance
func NewInMemoryStorage() *InMemoryStorage {
	source := rand.NewSource(time.Now().UnixNano())
	return &InMemoryStorage{
		todos: make(map[ulid.ULID]*todov1.Todo),
		rnd:   rand.New(source),
	}
}

// Add creates a new todo with the given title
func (s *InMemoryStorage) Add(title string) (*todov1.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entropy := ulid.Monotonic(s.rnd, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	todo := &todov1.Todo{
		Id:        id.String(),
		Title:     title,
		Completed: false,
	}
	s.todos[id] = todo

	return todo, nil
}

// Get retrieves a todo by ID
func (s *InMemoryStorage) Get(id ulid.ULID) (*todov1.Todo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, exists := s.todos[id]
	return todo, exists
}

// List returns all todos sorted by ID
func (s *InMemoryStorage) List() ([]*todov1.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a slice of all todos
	todos := make([]*todov1.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}

	// Sort by ULID
	sort.Slice(todos, func(i, j int) bool {
		idI, _ := ulid.Parse(todos[i].Id)
		idJ, _ := ulid.Parse(todos[j].Id)
		return idI.Compare(idJ) < 0
	})

	return todos, nil
}

// Update modifies a todo's title
func (s *InMemoryStorage) Update(id ulid.ULID, title string) (bool, error) {
	if title == "" {
		return false, fmt.Errorf("title cannot be empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return false, nil
	}

	todo.Title = title
	return true, nil
}

// Delete removes a todo
func (s *InMemoryStorage) Delete(id ulid.ULID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		return false, nil
	}

	delete(s.todos, id)
	return true, nil
}

// Complete marks a todo as completed
func (s *InMemoryStorage) Complete(id ulid.ULID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[id]
	if !exists {
		return false, nil
	}

	todo.Completed = true
	return true, nil
}
