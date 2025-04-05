package storage

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

// BenchmarkInMemoryStorage_Add benchmarks the Add method
func BenchmarkInMemoryStorage_Add(b *testing.B) {
	storage := NewInMemoryStorage()

	// Reset timer before starting the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Add a new todo with a unique title to avoid any caching
		title := "Todo" + ulid.MustNew(uint64(i), nil).String()
		_, err := storage.Add(title)
		if err != nil {
			b.Fatalf("Add failed: %v", err)
		}
	}
}

// BenchmarkInMemoryStorage_Get benchmarks the Get method
func BenchmarkInMemoryStorage_Get(b *testing.B) {
	storage := NewInMemoryStorage()

	// Add 100 todos to retrieve from
	ids := make([]ulid.ULID, 100)
	for i := 0; i < 100; i++ {
		todo, _ := storage.Add("Todo" + ulid.MustNew(uint64(i), nil).String())
		id, _ := ulid.Parse(todo.Id)
		ids[i] = id
	}

	// Reset timer before starting the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Get a random todo from the set
		id := ids[i%100]
		_, exists := storage.Get(id)
		if !exists {
			b.Fatalf("Get failed for ID: %s", id.String())
		}
	}
}

// BenchmarkInMemoryStorage_List benchmarks the List method
func BenchmarkInMemoryStorage_List(b *testing.B) {
	storage := NewInMemoryStorage()

	// Add 1000 todos to list
	for i := 0; i < 1000; i++ {
		_, _ = storage.Add("Todo" + ulid.MustNew(uint64(i), nil).String())
	}

	// Reset timer before starting the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		todos, err := storage.List()
		if err != nil {
			b.Fatalf("List failed: %v", err)
		}
		if len(todos) != 1000 {
			b.Fatalf("Expected 1000 todos, got %d", len(todos))
		}
	}
}

// BenchmarkInMemoryStorage_Update benchmarks the Update method
func BenchmarkInMemoryStorage_Update(b *testing.B) {
	storage := NewInMemoryStorage()

	// Add 100 todos to update
	ids := make([]ulid.ULID, 100)
	for i := 0; i < 100; i++ {
		todo, _ := storage.Add("Todo" + ulid.MustNew(uint64(i), nil).String())
		id, _ := ulid.Parse(todo.Id)
		ids[i] = id
	}

	// Reset timer before starting the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Update a random todo
		id := ids[i%100]
		newTitle := "Updated" + ulid.MustNew(uint64(i), nil).String()
		success, err := storage.Update(id, newTitle)
		if err != nil {
			b.Fatalf("Update failed: %v", err)
		}
		if !success {
			b.Fatalf("Update unsuccessful for ID: %s", id.String())
		}
	}
}

// BenchmarkInMemoryStorage_Delete benchmarks the Delete method
func BenchmarkInMemoryStorage_Delete(b *testing.B) {
	storage := NewInMemoryStorage()

	// Add 100 todos to delete
	ids := make([]ulid.ULID, 100)
	for i := 0; i < 100; i++ {
		todo, _ := storage.Add("Todo" + ulid.MustNew(uint64(i), nil).String())
		id, _ := ulid.Parse(todo.Id)
		ids[i] = id
	}

	// Reset timer before starting the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Delete a todo (cycling through the 100 we created)
		id := ids[i%100]

		// Clone the ID for use in deletion to avoid modifying the original array
		cloneID := ulid.MustParse(id.String())

		_, _ = storage.Delete(cloneID)
	}
}

// BenchmarkInMemoryStorage_Complete benchmarks the Complete method
func BenchmarkInMemoryStorage_Complete(b *testing.B) {
	storage := NewInMemoryStorage()

	// Add 100 todos to complete
	ids := make([]ulid.ULID, 100)
	for i := 0; i < 100; i++ {
		todo, _ := storage.Add("Todo" + ulid.MustNew(uint64(i), nil).String())
		id, _ := ulid.Parse(todo.Id)
		ids[i] = id
	}

	// Reset timer before starting the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Complete a random todo
		id := ids[i%100]
		success, err := storage.Complete(id)
		if err != nil {
			b.Fatalf("Complete failed: %v", err)
		}
		if !success {
			b.Fatalf("Complete unsuccessful for ID: %s", id.String())
		}
	}
}

// BenchmarkInMemoryStorage_Parallel benchmarks concurrent operations
func BenchmarkInMemoryStorage_Parallel(b *testing.B) {
	storage := NewInMemoryStorage()

	// Add some initial todos
	ids := make([]ulid.ULID, 100)
	for i := 0; i < 100; i++ {
		todo, _ := storage.Add("Todo" + ulid.MustNew(uint64(i), nil).String())
		id, _ := ulid.Parse(todo.Id)
		ids[i] = id
	}

	// Reset timer before starting the benchmark
	b.ResetTimer()

	// Run parallel benchmark
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			i++
			switch i % 5 {
			case 0:
				// Add
				_, _ = storage.Add("New Todo " + ulid.MustNew(ulid.Now(), nil).String())
			case 1:
				// Get
				id := ids[i%100]
				_, _ = storage.Get(id)
			case 2:
				// Update
				id := ids[i%100]
				_, _ = storage.Update(id, "Updated"+ulid.MustNew(ulid.Now(), nil).String())
			case 3:
				// List
				_, _ = storage.List()
			case 4:
				// Complete
				id := ids[i%100]
				_, _ = storage.Complete(id)
			}
		}
	})
}
