package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/oklog/ulid/v2"
	todov1 "github.com/scrogson/todo-golang/pkg/todo/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the TodoStorage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Add(title string) (*todov1.Todo, error) {
	args := m.Called(title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*todov1.Todo), args.Error(1)
}

func (m *MockStorage) Get(id ulid.ULID) (*todov1.Todo, bool) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*todov1.Todo), args.Bool(1)
}

func (m *MockStorage) List() ([]*todov1.Todo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*todov1.Todo), args.Error(1)
}

func (m *MockStorage) Update(id ulid.ULID, title string) (bool, error) {
	args := m.Called(id, title)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) Delete(id ulid.ULID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) Complete(id ulid.ULID) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestListTodos(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Create mock storage
		mockStorage := new(MockStorage)

		// Create sample todos
		todos := []*todov1.Todo{
			{Id: "01HZFG1EAQK0VKPNKN5AHF3QKP", Title: "Test Todo 1", Completed: false},
			{Id: "01HZFG1EAQK0VKPNKN5AHF3QKQ", Title: "Test Todo 2", Completed: true},
		}

		// Setup expectation
		mockStorage.On("List").Return(todos, nil)

		// Create server with mock storage
		server := NewTodoServer(mockStorage)

		// Call the method
		resp, err := server.ListTodos(context.Background(), &todov1.ListTodosRequest{})

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, todos, resp.Todos)

		// Verify expectations were met
		mockStorage.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Create mock storage
		mockStorage := new(MockStorage)

		// Setup expectation for error case
		expectedErr := fmt.Errorf("database connection error")
		mockStorage.On("List").Return(nil, expectedErr)

		// Create server with mock storage
		server := NewTodoServer(mockStorage)

		// Call the method
		resp, err := server.ListTodos(context.Background(), &todov1.ListTodosRequest{})

		// Verify results
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)

		// Verify expectations were met
		mockStorage.AssertExpectations(t)
	})
}

func TestAddTodo(t *testing.T) {
	// Create mock storage
	mockStorage := new(MockStorage)

	// Setup test cases
	testCases := []struct {
		name      string
		title     string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "Valid title",
			title: "New Todo",
			mockSetup: func() {
				mockStorage.On("Add", "New Todo").Return(&todov1.Todo{
					Id:        "01HZFG1EAQK0VKPNKN5AHF3QKR",
					Title:     "New Todo",
					Completed: false,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "Empty title",
			title:     "",
			mockSetup: func() {},
			wantErr:   true,
		},
		{
			name:  "Storage error",
			title: "Error Todo",
			mockSetup: func() {
				mockStorage.On("Add", "Error Todo").Return(nil, fmt.Errorf("storage error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock
			mockStorage = new(MockStorage)

			// Setup mock expectation
			tc.mockSetup()

			// Create server with mock storage
			server := NewTodoServer(mockStorage)

			// Call the method
			resp, err := server.AddTodo(context.Background(), &todov1.AddTodoRequest{
				Title: tc.title,
			})

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.title, resp.Todo.Title)
			}

			// Verify expectations were met
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	// Create mock storage
	mockStorage := new(MockStorage)

	// Valid ULID
	validID := "01HZFG1EAQK0VKPNKN5AHF3QKR"
	parsedID, _ := ulid.Parse(validID)

	// Setup test cases
	testCases := []struct {
		name      string
		id        string
		mockSetup func()
		wantErr   bool
		success   bool
	}{
		{
			name: "Valid delete",
			id:   validID,
			mockSetup: func() {
				mockStorage.On("Delete", parsedID).Return(true, nil)
			},
			wantErr: false,
			success: true,
		},
		{
			name:      "Invalid ID",
			id:        "invalid-id",
			mockSetup: func() {},
			wantErr:   true,
			success:   false,
		},
		{
			name: "Not found",
			id:   validID,
			mockSetup: func() {
				mockStorage.On("Delete", parsedID).Return(false, nil)
			},
			wantErr: false,
			success: false,
		},
		{
			name: "Storage error",
			id:   validID,
			mockSetup: func() {
				mockStorage.On("Delete", parsedID).Return(false, fmt.Errorf("storage error"))
			},
			wantErr: true,
			success: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock
			mockStorage = new(MockStorage)

			// Setup mock expectation
			tc.mockSetup()

			// Create server with mock storage
			server := NewTodoServer(mockStorage)

			// Call the method
			resp, err := server.DeleteTodo(context.Background(), &todov1.DeleteTodoRequest{
				Id: tc.id,
			})

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.success, resp.Success)
			}

			// Verify expectations were met
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	// Create mock storage
	mockStorage := new(MockStorage)

	// Valid ULID
	validID := "01HZFG1EAQK0VKPNKN5AHF3QKR"
	parsedID, _ := ulid.Parse(validID)

	// Setup test cases
	testCases := []struct {
		name      string
		id        string
		title     string
		mockSetup func()
		wantErr   bool
		success   bool
	}{
		{
			name:  "Valid update",
			id:    validID,
			title: "Updated Todo",
			mockSetup: func() {
				mockStorage.On("Update", parsedID, "Updated Todo").Return(true, nil)
			},
			wantErr: false,
			success: true,
		},
		{
			name:      "Invalid ID",
			id:        "invalid-id",
			title:     "Updated Todo",
			mockSetup: func() {},
			wantErr:   true,
			success:   false,
		},
		{
			name:  "Not found",
			id:    validID,
			title: "Updated Todo",
			mockSetup: func() {
				mockStorage.On("Update", parsedID, "Updated Todo").Return(false, nil)
			},
			wantErr: false,
			success: false,
		},
		{
			name:  "Storage error",
			id:    validID,
			title: "Error Todo",
			mockSetup: func() {
				mockStorage.On("Update", parsedID, "Error Todo").Return(false, fmt.Errorf("storage error"))
			},
			wantErr: true,
			success: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock
			mockStorage = new(MockStorage)

			// Setup mock expectation
			tc.mockSetup()

			// Create server with mock storage
			server := NewTodoServer(mockStorage)

			// Call the method
			resp, err := server.UpdateTodo(context.Background(), &todov1.UpdateTodoRequest{
				Id:    tc.id,
				Title: tc.title,
			})

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.success, resp.Success)
			}

			// Verify expectations were met
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestCompleteTodo(t *testing.T) {
	// Create mock storage
	mockStorage := new(MockStorage)

	// Valid ULID
	validID := "01HZFG1EAQK0VKPNKN5AHF3QKR"
	parsedID, _ := ulid.Parse(validID)

	// Setup test cases
	testCases := []struct {
		name      string
		id        string
		mockSetup func()
		wantErr   bool
		success   bool
	}{
		{
			name: "Valid complete",
			id:   validID,
			mockSetup: func() {
				mockStorage.On("Complete", parsedID).Return(true, nil)
			},
			wantErr: false,
			success: true,
		},
		{
			name:      "Invalid ID",
			id:        "invalid-id",
			mockSetup: func() {},
			wantErr:   true,
			success:   false,
		},
		{
			name: "Not found",
			id:   validID,
			mockSetup: func() {
				mockStorage.On("Complete", parsedID).Return(false, nil)
			},
			wantErr: false,
			success: false,
		},
		{
			name: "Storage error",
			id:   validID,
			mockSetup: func() {
				mockStorage.On("Complete", parsedID).Return(false, fmt.Errorf("storage error"))
			},
			wantErr: true,
			success: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock
			mockStorage = new(MockStorage)

			// Setup mock expectation
			tc.mockSetup()

			// Create server with mock storage
			server := NewTodoServer(mockStorage)

			// Call the method
			resp, err := server.CompleteTodo(context.Background(), &todov1.CompleteTodoRequest{
				Id: tc.id,
			})

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.success, resp.Success)
			}

			// Verify expectations were met
			mockStorage.AssertExpectations(t)
		})
	}
}
