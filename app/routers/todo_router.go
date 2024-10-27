package routers

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

type TodoRouter interface {
	SetRouting(r *echo.Echo)
}

type todoRouter struct {
	todoController controllers.TodoController
}

func NewTodoRouter(todoController controllers.TodoController) TodoRouter {
	return &todoRouter{todoController}
}

func (tr *todoRouter) SetRouting(r *echo.Echo) {
	g := r.Group("/todos")
	g.POST("", tr.todoController.Create)
	g.GET("", tr.todoController.Index)
	g.GET("/:id", tr.todoController.Show)
	g.PUT("/:id", tr.todoController.Update)
	g.DELETE("/:id", tr.todoController.Delete)
}
