package server

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/scrogson/todo-golang/internal/storage"
	todov1 "github.com/scrogson/todo-golang/pkg/todo/v1"
)

// TodoServer implements the TodoService gRPC service
type TodoServer struct {
	todov1.UnimplementedTodoServiceServer
	storage storage.TodoStorage
}

// NewTodoServer creates a new TodoServer
func NewTodoServer(storage storage.TodoStorage) *TodoServer {
	return &TodoServer{
		storage: storage,
	}
}

// ListTodos returns all todos
func (s *TodoServer) ListTodos(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	resp := &todov1.ListTodosResponse{}

	todos, err := s.storage.List()
	if err != nil {
		return nil, err
	}

	resp.Todos = append(resp.Todos, todos...)

	return resp, nil
}

// AddTodo creates a new todo
func (s *TodoServer) AddTodo(ctx context.Context, req *todov1.AddTodoRequest) (*todov1.AddTodoResponse, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title cannot be empty")
	}

	todo, err := s.storage.Add(req.Title)
	if err != nil {
		return nil, err
	}

	return &todov1.AddTodoResponse{Todo: todo}, nil
}

// DeleteTodo removes a todo by ID
func (s *TodoServer) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	id, err := ulid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %s", err)
	}

	success, err := s.storage.Delete(id)
	if err != nil {
		return nil, err
	}

	return &todov1.DeleteTodoResponse{Success: success}, nil
}

// UpdateTodo updates a todo's title
func (s *TodoServer) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	id, err := ulid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %s", err)
	}

	success, err := s.storage.Update(id, req.Title)
	if err != nil {
		return nil, err
	}

	return &todov1.UpdateTodoResponse{Success: success}, nil
}

// CompleteTodo marks a todo as completed
func (s *TodoServer) CompleteTodo(ctx context.Context, req *todov1.CompleteTodoRequest) (*todov1.CompleteTodoResponse, error) {
	id, err := ulid.Parse(req.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %s", err)
	}

	success, err := s.storage.Complete(id)
	if err != nil {
		return nil, err
	}

	return &todov1.CompleteTodoResponse{Success: success}, nil
}
