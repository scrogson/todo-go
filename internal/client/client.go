package client

import (
	"context"

	todov1 "github.com/scrogson/todo-go/pkg/todo/v1"
)

// TodoClient wraps the gRPC client with helpful methods
type TodoClient struct {
	client todov1.TodoServiceClient
}

// NewTodoClient creates a new TodoClient
func NewTodoClient(client todov1.TodoServiceClient) *TodoClient {
	return &TodoClient{client: client}
}

// ListTodos fetches all todos
func (c *TodoClient) ListTodos(ctx context.Context) ([]*todov1.Todo, error) {
	resp, err := c.client.ListTodos(ctx, &todov1.ListTodosRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Todos, nil
}

// AddTodo creates a new todo
func (c *TodoClient) AddTodo(ctx context.Context, title string) (*todov1.Todo, error) {
	resp, err := c.client.AddTodo(ctx, &todov1.AddTodoRequest{Title: title})
	if err != nil {
		return nil, err
	}
	return resp.Todo, nil
}

// DeleteTodo deletes a todo by ID
func (c *TodoClient) DeleteTodo(ctx context.Context, id string) (bool, error) {
	resp, err := c.client.DeleteTodo(ctx, &todov1.DeleteTodoRequest{Id: id})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

// UpdateTodo updates a todo's title
func (c *TodoClient) UpdateTodo(ctx context.Context, id, title string) (bool, error) {
	resp, err := c.client.UpdateTodo(ctx, &todov1.UpdateTodoRequest{Id: id, Title: title})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

// CompleteTodo marks a todo as complete
func (c *TodoClient) CompleteTodo(ctx context.Context, id string) (bool, error) {
	resp, err := c.client.CompleteTodo(ctx, &todov1.CompleteTodoRequest{Id: id})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
