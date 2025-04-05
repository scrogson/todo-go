package client

import (
	"context"
	"errors"
	"testing"

	todov1 "github.com/scrogson/todo-golang/pkg/todo/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockTodoServiceClient is a mock implementation of TodoServiceClient
type MockTodoServiceClient struct {
	mock.Mock
}

func (m *MockTodoServiceClient) ListTodos(ctx context.Context, req *todov1.ListTodosRequest, opts ...grpc.CallOption) (*todov1.ListTodosResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*todov1.ListTodosResponse), args.Error(1)
}

func (m *MockTodoServiceClient) AddTodo(ctx context.Context, req *todov1.AddTodoRequest, opts ...grpc.CallOption) (*todov1.AddTodoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*todov1.AddTodoResponse), args.Error(1)
}

func (m *MockTodoServiceClient) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest, opts ...grpc.CallOption) (*todov1.DeleteTodoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*todov1.DeleteTodoResponse), args.Error(1)
}

func (m *MockTodoServiceClient) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest, opts ...grpc.CallOption) (*todov1.UpdateTodoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*todov1.UpdateTodoResponse), args.Error(1)
}

func (m *MockTodoServiceClient) CompleteTodo(ctx context.Context, req *todov1.CompleteTodoRequest, opts ...grpc.CallOption) (*todov1.CompleteTodoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*todov1.CompleteTodoResponse), args.Error(1)
}

func TestListTodos(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	todoClient := NewTodoClient(mockClient)
	ctx := context.Background()

	// Test successful response
	todos := []*todov1.Todo{
		{Id: "todo1", Title: "First Todo", Completed: false},
		{Id: "todo2", Title: "Second Todo", Completed: true},
	}
	mockClient.On("ListTodos", ctx, &todov1.ListTodosRequest{}).Return(&todov1.ListTodosResponse{
		Todos: todos,
	}, nil)

	result, err := todoClient.ListTodos(ctx)
	assert.NoError(t, err)
	assert.Equal(t, todos, result)

	// Test error response
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	expectedErr := errors.New("connection error")
	mockClient.On("ListTodos", ctx, &todov1.ListTodosRequest{}).Return(nil, expectedErr)

	result, err = todoClient.ListTodos(ctx)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)

	mockClient.AssertExpectations(t)
}

func TestAddTodo(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	todoClient := NewTodoClient(mockClient)
	ctx := context.Background()
	title := "New Todo"

	// Test successful response
	todo := &todov1.Todo{Id: "new-id", Title: title, Completed: false}
	mockClient.On("AddTodo", ctx, &todov1.AddTodoRequest{Title: title}).Return(&todov1.AddTodoResponse{
		Todo: todo,
	}, nil)

	result, err := todoClient.AddTodo(ctx, title)
	assert.NoError(t, err)
	assert.Equal(t, todo, result)

	// Test error response
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	expectedErr := errors.New("connection error")
	mockClient.On("AddTodo", ctx, &todov1.AddTodoRequest{Title: title}).Return(nil, expectedErr)

	result, err = todoClient.AddTodo(ctx, title)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)

	mockClient.AssertExpectations(t)
}

func TestDeleteTodo(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	todoClient := NewTodoClient(mockClient)
	ctx := context.Background()
	id := "todo-id"

	// Test successful response - deleted
	mockClient.On("DeleteTodo", ctx, &todov1.DeleteTodoRequest{Id: id}).Return(&todov1.DeleteTodoResponse{
		Success: true,
	}, nil)

	success, err := todoClient.DeleteTodo(ctx, id)
	assert.NoError(t, err)
	assert.True(t, success)

	// Test successful response - not found
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	mockClient.On("DeleteTodo", ctx, &todov1.DeleteTodoRequest{Id: id}).Return(&todov1.DeleteTodoResponse{
		Success: false,
	}, nil)

	success, err = todoClient.DeleteTodo(ctx, id)
	assert.NoError(t, err)
	assert.False(t, success)

	// Test error response
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	expectedErr := errors.New("connection error")
	mockClient.On("DeleteTodo", ctx, &todov1.DeleteTodoRequest{Id: id}).Return(nil, expectedErr)

	success, err = todoClient.DeleteTodo(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.False(t, success)

	mockClient.AssertExpectations(t)
}

func TestUpdateTodo(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	todoClient := NewTodoClient(mockClient)
	ctx := context.Background()
	id := "todo-id"
	title := "Updated Title"

	// Test successful response - updated
	mockClient.On("UpdateTodo", ctx, &todov1.UpdateTodoRequest{Id: id, Title: title}).Return(&todov1.UpdateTodoResponse{
		Success: true,
	}, nil)

	success, err := todoClient.UpdateTodo(ctx, id, title)
	assert.NoError(t, err)
	assert.True(t, success)

	// Test successful response - not found
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	mockClient.On("UpdateTodo", ctx, &todov1.UpdateTodoRequest{Id: id, Title: title}).Return(&todov1.UpdateTodoResponse{
		Success: false,
	}, nil)

	success, err = todoClient.UpdateTodo(ctx, id, title)
	assert.NoError(t, err)
	assert.False(t, success)

	// Test error response
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	expectedErr := errors.New("connection error")
	mockClient.On("UpdateTodo", ctx, &todov1.UpdateTodoRequest{Id: id, Title: title}).Return(nil, expectedErr)

	success, err = todoClient.UpdateTodo(ctx, id, title)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.False(t, success)

	mockClient.AssertExpectations(t)
}

func TestCompleteTodo(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	todoClient := NewTodoClient(mockClient)
	ctx := context.Background()
	id := "todo-id"

	// Test successful response - completed
	mockClient.On("CompleteTodo", ctx, &todov1.CompleteTodoRequest{Id: id}).Return(&todov1.CompleteTodoResponse{
		Success: true,
	}, nil)

	success, err := todoClient.CompleteTodo(ctx, id)
	assert.NoError(t, err)
	assert.True(t, success)

	// Test successful response - not found
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	mockClient.On("CompleteTodo", ctx, &todov1.CompleteTodoRequest{Id: id}).Return(&todov1.CompleteTodoResponse{
		Success: false,
	}, nil)

	success, err = todoClient.CompleteTodo(ctx, id)
	assert.NoError(t, err)
	assert.False(t, success)

	// Test error response
	mockClient = new(MockTodoServiceClient)
	todoClient = NewTodoClient(mockClient)
	expectedErr := errors.New("connection error")
	mockClient.On("CompleteTodo", ctx, &todov1.CompleteTodoRequest{Id: id}).Return(nil, expectedErr)

	success, err = todoClient.CompleteTodo(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.False(t, success)

	mockClient.AssertExpectations(t)
}
