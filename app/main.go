package main

import (
	"app/controllers"
	"app/db"
	"app/routers"
	"app/services"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	dbCon := db.Init()

	// service
	authService := services.NewAuthService(dbCon)
	todoService := services.NewTodoService(dbCon)

	// controller
	authController := controllers.NewAuthController(authService)
	todoController := controllers.NewTodoController(todoService, authService)

	// router
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	authRouter := routers.NewAuthRouter(authController)
	todoRouter := routers.NewTodoRouter(todoController)
	authRouter.SetRouting(e)
	todoRouter.SetRouting(e)

	e.Logger.Fatal(e.Start(":" + os.Getenv("SERVER_PORT")))
}
