package storage

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	todov1 "github.com/scrogson/todo-go/pkg/todo/v1"
)

// SQLiteStorage implements TodoStorage interface with SQLite database
type SQLiteStorage struct {
	db  *sql.DB
	rnd *rand.Rand
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create todos table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create todos table: %w", err)
	}

	source := rand.NewSource(time.Now().UnixNano())
	return &SQLiteStorage{
		db:  db,
		rnd: rand.New(source),
	}, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// Add creates a new todo with the given title
func (s *SQLiteStorage) Add(title string) (*todov1.Todo, error) {
	entropy := ulid.Monotonic(s.rnd, 0)
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	todo := &todov1.Todo{
		Id:        id.String(),
		Title:     title,
		Completed: false,
	}

	_, err := s.db.Exec("INSERT INTO todos (id, title, completed) VALUES (?, ?, ?)",
		todo.Id, todo.Title, todo.Completed)
	if err != nil {
		return nil, fmt.Errorf("failed to add todo: %w", err)
	}

	return todo, nil
}

// Get retrieves a todo by ID
func (s *SQLiteStorage) Get(id ulid.ULID) (*todov1.Todo, bool) {
	var todo todov1.Todo
	err := s.db.QueryRow("SELECT id, title, completed FROM todos WHERE id = ?", id.String()).
		Scan(&todo.Id, &todo.Title, &todo.Completed)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		// Log the error but return false to conform to the interface
		return nil, false
	}

	return &todo, true
}

// List returns all todos sorted by ID
func (s *SQLiteStorage) List() ([]*todov1.Todo, error) {
	rows, err := s.db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		return nil, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []*todov1.Todo
	for rows.Next() {
		var todo todov1.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Completed); err != nil {
			return nil, fmt.Errorf("failed to scan todo row: %w", err)
		}
		todos = append(todos, &todo)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
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
func (s *SQLiteStorage) Update(id ulid.ULID, title string) (bool, error) {
	if title == "" {
		return false, fmt.Errorf("title cannot be empty")
	}

	result, err := s.db.Exec("UPDATE todos SET title = ? WHERE id = ?", title, id.String())
	if err != nil {
		return false, fmt.Errorf("failed to update todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}

// Delete removes a todo
func (s *SQLiteStorage) Delete(id ulid.ULID) (bool, error) {
	result, err := s.db.Exec("DELETE FROM todos WHERE id = ?", id.String())
	if err != nil {
		return false, fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}

// Complete marks a todo as completed
func (s *SQLiteStorage) Complete(id ulid.ULID) (bool, error) {
	result, err := s.db.Exec("UPDATE todos SET completed = 1 WHERE id = ?", id.String())
	if err != nil {
		return false, fmt.Errorf("failed to complete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
