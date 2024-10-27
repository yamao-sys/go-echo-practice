package routers

import (
	"app/controllers"

	"github.com/labstack/echo/v4"
)

type AuthRouter interface {
	SetRouting(r *echo.Echo)
}

type authRouter struct {
	authController controllers.AuthController
}

func NewAuthRouter(authController controllers.AuthController) AuthRouter {
	return &authRouter{authController}
}

func (ar *authRouter) SetRouting(r *echo.Echo) {
	g := r.Group("/auth")
	g.POST("/sign_up", ar.authController.SignUp)
	g.POST("/sign_in", ar.authController.SignIn)
}
