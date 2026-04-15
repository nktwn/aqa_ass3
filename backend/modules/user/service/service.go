package user

import (
	"context"
	"diploma/modules/user/model"
	"diploma/pkg/client/db"
)

type userServ struct {
	userRepository IUserRepository
	txManager      db.TxManager
}

func NewService(
	userRepository IUserRepository,
	txManager db.TxManager,
) *userServ {
	return &userServ{
		userRepository: userRepository,
		txManager:      txManager,
	}
}

type IUserRepository interface {
	UserByID(ctx context.Context, userID int64) (model.User, error)
	UpdateUser(ctx context.Context, user model.User)  error
	IAddressRepo
}
