package dto

import models "app/models/generated"

type SignUpRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUpResponse struct {
	User      models.User
	Error     error
	ErrorType string
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	TokenString     string
	NotFoundMessage string
	Error           error
}
