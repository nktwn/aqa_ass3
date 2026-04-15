package auth

import (
	"context"
	"diploma/modules/auth/model"
	"diploma/pkg/client/db"
	"diploma/pkg/service"
)

type authServ struct {
	service.BaseService
	authRepository IAuthRepository
	jwt            IJWT
	txManager      db.TxManager
}

func NewService(
	authRepository IAuthRepository,
	jwt IJWT,
	txManager db.TxManager,
) *authServ {
	return &authServ{
		BaseService:    service.NewBaseService("auth"),
		authRepository: authRepository,
		jwt:            jwt,
		txManager:      txManager,
	}
}

type IAuthRepository interface {
	Create(ctx context.Context, user *model.AuthUser) (int64, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.AuthUser, error)
}

type IJWT interface {
	GenerateJSONWebTokens(id int64, username string, role int) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(refreshToken string) (string, error)
	// VerifyToken(accessToken string) (bool, error)
}
