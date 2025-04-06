package storage

import (
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
)

func BenchmarkSQLiteStorage_Add(b *testing.B) {
	// Use an in-memory database for benchmarking
	storage, err := NewSQLiteStorage(":memory:")
	if err != nil {
		b.Fatalf("Failed to create SQLite storage: %v", err)
	}
	defer storage.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.Add("Test Todo")
		if err != nil {
			b.Fatalf("Failed to add todo: %v", err)
		}
	}
}

func BenchmarkSQLiteStorage_Get(b *testing.B) {
	// Use an in-memory database for benchmarking
	storage, err := NewSQLiteStorage(":memory:")
	if err != nil {
		b.Fatalf("Failed to create SQLite storage: %v", err)
	}
	defer storage.Close()

	// Add a todo to get
	todo, err := storage.Add("Test Todo")
	if err != nil {
		b.Fatalf("Failed to add todo: %v", err)
	}
	id, err := ulid.Parse(todo.Id)
	if err != nil {
		b.Fatalf("Failed to parse ULID: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, exists := storage.Get(id)
		if !exists {
			b.Fatalf("Todo not found")
		}
	}
}

func BenchmarkSQLiteStorage_List(b *testing.B) {
	// Use an in-memory database for benchmarking
	storage, err := NewSQLiteStorage(":memory:")
	if err != nil {
		b.Fatalf("Failed to create SQLite storage: %v", err)
	}
	defer storage.Close()

	// Add some todos
	for i := 0; i < 50; i++ {
		_, err := storage.Add("Test Todo")
		if err != nil {
			b.Fatalf("Failed to add todo: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := storage.List()
		if err != nil {
			b.Fatalf("Failed to list todos: %v", err)
		}
	}
}

func BenchmarkSQLiteStorage_Update(b *testing.B) {
	// Use an in-memory database for benchmarking
	storage, err := NewSQLiteStorage(":memory:")
	if err != nil {
		b.Fatalf("Failed to create SQLite storage: %v", err)
	}
	defer storage.Close()

	// Add a todo to update
	todo, err := storage.Add("Test Todo")
	if err != nil {
		b.Fatalf("Failed to add todo: %v", err)
	}
	id, err := ulid.Parse(todo.Id)
	if err != nil {
		b.Fatalf("Failed to parse ULID: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updated, err := storage.Update(id, "Updated Todo")
		if err != nil {
			b.Fatalf("Failed to update todo: %v", err)
		}
		if !updated {
			b.Fatalf("Todo not updated")
		}
	}
}

func BenchmarkSQLiteStorage_Delete(b *testing.B) {
	b.StopTimer()

	tempFile, err := os.CreateTemp("", "todo-sqlite-bench-*.db")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()
	defer os.Remove(dbPath)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create a new storage for each iteration to avoid interference
		storage, err := NewSQLiteStorage(dbPath)
		if err != nil {
			b.Fatalf("Failed to create SQLite storage: %v", err)
		}

		// Add a todo to delete
		todo, err := storage.Add("Test Todo")
		if err != nil {
			b.Fatalf("Failed to add todo: %v", err)
		}
		id, err := ulid.Parse(todo.Id)
		if err != nil {
			b.Fatalf("Failed to parse ULID: %v", err)
		}

		b.StartTimer()
		_, err = storage.Delete(id)
		if err != nil {
			b.Fatalf("Failed to delete todo: %v", err)
		}
		b.StopTimer()

		storage.Close()
	}
}

func BenchmarkSQLiteStorage_Complete(b *testing.B) {
	// Use an in-memory database for benchmarking
	storage, err := NewSQLiteStorage(":memory:")
	if err != nil {
		b.Fatalf("Failed to create SQLite storage: %v", err)
	}
	defer storage.Close()

	// Add a todo to complete
	todo, err := storage.Add("Test Todo")
	if err != nil {
		b.Fatalf("Failed to add todo: %v", err)
	}
	id, err := ulid.Parse(todo.Id)
	if err != nil {
		b.Fatalf("Failed to parse ULID: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Reset the todo for each iteration
		_, err := storage.Update(id, "Test Todo")
		if err != nil {
			b.Fatalf("Failed to reset todo: %v", err)
		}
		b.StartTimer()

		completed, err := storage.Complete(id)
		if err != nil {
			b.Fatalf("Failed to complete todo: %v", err)
		}
		if !completed {
			b.Fatalf("Todo not completed")
		}
	}
}
