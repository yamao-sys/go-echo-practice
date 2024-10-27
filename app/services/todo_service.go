package services

import (
	"app/dto"
	models "app/models/generated"
	"context"
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type TodoService interface {
	CreateTodo(ctx context.Context, requestParams dto.CreateTodoRequest, userID int) *dto.CreateTodoResponse
	FetchTodosList(ctx context.Context, userID int) *dto.TodosListResponse
	FetchTodo(ctx context.Context, id int, userID int) *dto.FetchTodoResponse
	UpdateTodo(ctx context.Context, id int, requestParams dto.UpdateTodoRequest, userID int) *dto.UpdateTodoResponse
	DeleteTodo(ctx context.Context, id int, userID int) *dto.DeleteTodoResponse
}

type todoService struct {
	db *sql.DB
}

func NewTodoService(db *sql.DB) TodoService {
	return &todoService{db}
}

func (ts *todoService) CreateTodo(ctx context.Context, requestParams dto.CreateTodoRequest, userID int) *dto.CreateTodoResponse {
	// NOTE: バリデーションチェック
	validate := validator.New()
	validationErrors := validate.Struct(requestParams)
	if validationErrors != nil {
		return &dto.CreateTodoResponse{Todo: &models.Todo{}, Error: validationErrors, ErrorType: "validationError"}
	}

	todo := &models.Todo{}
	todo.Title = requestParams.Title
	todo.Content = null.String{String: requestParams.Content, Valid: true}
	todo.UserID = userID
	// NOTE: Create処理
	err := todo.Insert(ctx, ts.db, boil.Infer())
	if err != nil {
		return &dto.CreateTodoResponse{Todo: todo, Error: err, ErrorType: "internalServerError"}
	}
	return &dto.CreateTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) FetchTodosList(ctx context.Context, userID int) *dto.TodosListResponse {
	todos, error := models.Todos(qm.Where("user_id = ?", userID)).All(ctx, ts.db)
	if error != nil {
		return &dto.TodosListResponse{Todos: models.TodoSlice{}, Error: error, ErrorType: "notFound"}
	}
	fmt.Printf("todos %v", todos)

	return &dto.TodosListResponse{Todos: todos, Error: nil, ErrorType: ""}
}

func (ts *todoService) FetchTodo(ctx context.Context, id int, userID int) *dto.FetchTodoResponse {
	todo, error := models.Todos(qm.Where("id = ? AND user_id = ?", id, userID)).One(ctx, ts.db)
	if error != nil {
		return &dto.FetchTodoResponse{Todo: &models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	return &dto.FetchTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) UpdateTodo(ctx context.Context, id int, requestParams dto.UpdateTodoRequest, userID int) *dto.UpdateTodoResponse {
	todo, error := models.Todos(qm.Where("id = ? AND user_id = ?", id, userID)).One(ctx, ts.db)
	if error != nil {
		return &dto.UpdateTodoResponse{Todo: &models.Todo{}, Error: error, ErrorType: "notFound"}
	}

	// NOTE: バリデーションチェック
	validate := validator.New()
	validationErrors := validate.Struct(requestParams)
	if validationErrors != nil {
		return &dto.UpdateTodoResponse{Todo: todo, Error: validationErrors, ErrorType: "validationError"}
	}

	todo.Title = requestParams.Title
	todo.Content = null.String{String: requestParams.Content, Valid: true}

	// NOTE: Update処理
	_, updateError := todo.Update(ctx, ts.db, boil.Infer())
	if updateError != nil {
		return &dto.UpdateTodoResponse{Todo: todo, Error: updateError, ErrorType: "internalServerError"}
	}
	return &dto.UpdateTodoResponse{Todo: todo, Error: nil, ErrorType: ""}
}

func (ts *todoService) DeleteTodo(ctx context.Context, id int, userID int) *dto.DeleteTodoResponse {
	todo, error := models.Todos(qm.Where("id = ? AND user_id = ?", id, userID)).One(ctx, ts.db)
	if error != nil {
		return &dto.DeleteTodoResponse{Error: error, ErrorType: "notFound"}
	}

	_, deleteError := todo.Delete(ctx, ts.db)
	if deleteError != nil {
		return &dto.DeleteTodoResponse{Error: deleteError, ErrorType: "internalServerError"}
	}
	return &dto.DeleteTodoResponse{Error: nil, ErrorType: ""}
}
