package controllers

import (
	"app/dto"
	"app/services"
	"app/utils"
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TodoController interface {
	Create(ctx echo.Context) error
	Index(ctx echo.Context) error
	Show(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type todoController struct {
	todoService services.TodoService
	authService services.AuthService
}

func NewTodoController(todoService services.TodoService, authService services.AuthService) TodoController {
	return &todoController{todoService, authService}
}

func (todoController *todoController) Create(ctx echo.Context) error {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, responseHash("unauthorized error"))
	}

	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.CreateTodoRequest{}
	if err := ctx.Bind(&requestParams); err != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash(err))
	}

	createTodoContext := context.Background()
	result := todoController.todoService.CreateTodo(createTodoContext, requestParams, user.ID)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, result.Todo)
	}

	switch result.ErrorType {
	case "internalServerError":
		return ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	case "validationError":
		return ctx.JSON(http.StatusBadRequest, responseHash(utils.CoordinateValidationErrors(result.Error)))
	}
	return ctx.JSON(http.StatusInternalServerError, responseHash("unexpected error"))
}

func (todoController *todoController) Index(ctx echo.Context) error {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, responseHash("unauthorized error"))
	}

	fetchTodosContext := context.Background()
	result := todoController.todoService.FetchTodosList(fetchTodosContext, user.ID)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, result.Todos)
	}

	switch result.ErrorType {
	case "notFound":
		return ctx.JSON(http.StatusNotFound, responseHash(result.Error))
	case "internalServerError":
		return ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	}
	return ctx.JSON(http.StatusInternalServerError, responseHash("unexpected error"))
}

func (todoController *todoController) Show(ctx echo.Context) error {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, responseHash("unauthorized error"))
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash("internal server error"))
	}
	showTodoContext := context.Background()
	result := todoController.todoService.FetchTodo(showTodoContext, id, user.ID)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, result.Todo)
	}

	switch result.ErrorType {
	case "notFound":
		return ctx.JSON(http.StatusNotFound, responseHash(result.Error))
	case "internalServerError":
		return ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	}
	return ctx.JSON(http.StatusInternalServerError, responseHash("unexpected error"))
}

func (todoController *todoController) Update(ctx echo.Context) error {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, responseHash("unauthorized error"))
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash("internal server error"))
	}
	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.UpdateTodoRequest{}
	if err := ctx.Bind(&requestParams); err != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash(err))
	}
	updateTodoContext := context.Background()
	result := todoController.todoService.UpdateTodo(updateTodoContext, id, requestParams, user.ID)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, result.Todo)
	}

	switch result.ErrorType {
	case "validationError":
		return ctx.JSON(http.StatusBadRequest, responseHash(utils.CoordinateValidationErrors(result.Error)))
	case "notFound":
		return ctx.JSON(http.StatusNotFound, responseHash(result.Error))
	case "internalServerError":
		return ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	}
	return ctx.JSON(http.StatusInternalServerError, responseHash("unexpected error"))
}

func (todoController *todoController) Delete(ctx echo.Context) error {
	user, err := todoController.authService.GetAuthUser(ctx)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, responseHash("unauthorized error"))
	}
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash("internal server error"))
	}
	deleteTodoContext := context.Background()
	result := todoController.todoService.DeleteTodo(deleteTodoContext, id, user.ID)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, responseHash(""))
	}

	switch result.ErrorType {
	case "notFound":
		ctx.JSON(http.StatusNotFound, responseHash(result.Error))
	case "internalServerError":
		ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	}
	return ctx.JSON(http.StatusInternalServerError, responseHash("unexpected error"))
}
