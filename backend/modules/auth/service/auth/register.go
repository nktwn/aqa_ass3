package auth

import (
	"context"
	"errors"

	"diploma/modules/auth/model"
	passwordutil "diploma/pkg/password-util"

	"go.uber.org/zap"
)

func (s *authServ) Register(ctx context.Context, user *model.AuthUser) (int64, error) {
	s.LogInfo(ctx, "Starting user registration", zap.String("phone_number", user.Info.PhoneNumber))

	var id int64
	var err error
	user.HashedPassword, err = passwordutil.HashPassword(user.Password)
	if err != nil {
		s.LogError(ctx, "Failed to hash password", err)
		return 0, err
	}

	if user.Info.Role == model.AdminRole {
		s.LogWarn(ctx, "Attempted to register with admin role", zap.String("phone_number", user.Info.PhoneNumber))
		return 0, errors.New("admin role is not allowed")
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.authRepository.Create(ctx, user)
		if errTx != nil {
			s.LogError(ctx, "Failed to create user", errTx, zap.String("phone_number", user.Info.PhoneNumber))
			return errTx
		}

		_, errTx = s.authRepository.GetById(ctx, id)
		if errTx != nil {
			s.LogError(ctx, "Failed to retrieve created user", errTx, zap.Int64("user_id", id))
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	s.LogInfo(ctx, "User registration successful",
		zap.Int64("user_id", id),
		zap.String("phone_number", user.Info.PhoneNumber),
		zap.String("name", user.Info.Name),
	)
	return id, nil
}

func (s *authServ) Login(ctx context.Context, phoneNumber string, password string) (accessToken string, refreshToken string, err error) {
	s.LogInfo(ctx, "Attempting user login", zap.String("phone_number", phoneNumber))

	authUser, err := s.authRepository.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, model.ErrNoRows) {
			s.LogWarn(ctx, "Login attempt with non-existent phone number", zap.String("phone_number", phoneNumber))
			return "", "", model.ErrInvalidCredentials
		}
		s.LogError(ctx, "Failed to retrieve user by phone number", err, zap.String("phone_number", phoneNumber))
		return "", "", err
	}

	if !passwordutil.CheckPasswordHash(password, authUser.HashedPassword) {
		s.LogWarn(ctx, "Invalid password attempt", zap.String("phone_number", phoneNumber))
		return "", "", model.ErrInvalidCredentials
	}

	accessToken, refreshToken, err = s.jwt.GenerateJSONWebTokens(authUser.ID, authUser.Info.Name, authUser.Info.Role)
	if err != nil {
		s.LogError(ctx, "Failed to generate JWT tokens", err, zap.Int64("user_id", authUser.ID))
		return "", "", err
	}

	s.LogInfo(ctx, "User login successful",
		zap.Int64("user_id", authUser.ID),
		zap.String("phone_number", phoneNumber),
		zap.String("name", authUser.Info.Name),
	)

	return accessToken, refreshToken, nil
}
