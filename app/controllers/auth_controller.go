package controllers

import (
	"app/dto"
	"app/services"
	"app/utils"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController interface {
	SignUp(ctx echo.Context) error
	SignIn(ctx echo.Context) error
}

type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{authService}
}

func (authController *authController) SignUp(ctx echo.Context) error {
	// NOTE: リクエストデータを構造体に変換
	requestParams := dto.SignUpRequest{}
	if err := ctx.Bind(&requestParams); err != nil {
		log.Println("mapping error")
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, responseHash(err))
	}
	singUpContext := context.Background()
	result := authController.authService.SignUp(singUpContext, requestParams)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, responseHash(""))
	}

	switch result.ErrorType {
	case "internalServerError":
		return ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	case "validationError":
		return ctx.JSON(http.StatusBadRequest, responseHash(utils.CoordinateValidationErrors(result.Error)))
	}
	return ctx.JSON(http.StatusInternalServerError, responseHash("unexpected error"))
}

func (authController *authController) SignIn(ctx echo.Context) error {
	requestParams := dto.SignInRequest{}
	if err := ctx.Bind(&requestParams); err != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash(err))
	}
	singInContext := context.Background()
	result := authController.authService.SignIn(singInContext, requestParams)

	if result.NotFoundMessage != "" {
		return ctx.JSON(http.StatusNotFound, responseHash(result.NotFoundMessage))
	}
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, responseHash(result.Error))
	}

	// NOTE: Cookieにtokenをセット
	cookie := &http.Cookie{
		Name:     "token",
		Value:    result.TokenString,
		MaxAge:   3600 * 24,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
	return ctx.JSON(http.StatusOK, responseHash(""))
}

func responseHash(error interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	result["error"] = error
	return result
}
