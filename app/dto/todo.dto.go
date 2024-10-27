package dto

import models "app/models/generated"

type CreateTodoRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type CreateTodoResponse struct {
	Todo      *models.Todo
	Error     error
	ErrorType string
}

type TodosListResponse struct {
	Todos     models.TodoSlice
	Error     error
	ErrorType string
}

type FetchTodoResponse struct {
	Todo      *models.Todo
	Error     error
	ErrorType string
}

type UpdateTodoRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type UpdateTodoResponse struct {
	Todo      *models.Todo
	Error     error
	ErrorType string
}

type DeleteTodoResponse struct {
	Error     error
	ErrorType string
}
