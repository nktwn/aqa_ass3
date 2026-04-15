package handler

import (
	"context"
	"diploma/modules/user/model"
)

type UserHandler struct {
	service IUserService
}

func NewHandler(service IUserService) *UserHandler {
	return &UserHandler{service: service}
}

type IUserService interface {
	User(ctx context.Context, userID int64) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
	IAddressService
	GetUserRole(ctx context.Context, userID int64) (int, error)
}
