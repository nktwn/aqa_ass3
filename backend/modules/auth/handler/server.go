package handler

import (
	"context"
	"diploma/modules/auth/model"
)

type AuthHandler struct {
	service IAuthService
}

func NewHandler(service IAuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type IAuthService interface {
	Register(ctx context.Context, user *model.AuthUser) (int64, error)
	Login(ctx context.Context, phoneNumber string, password string) (accessToken string, refreshToken string, err error)
}
