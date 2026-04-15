package user

import (
	"context"
	"diploma/modules/user/model"
)

func (s *userServ) User(ctx context.Context, userID int64) (model.User, error) {
	user, err := s.userRepository.UserByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	addressList, err := s.userRepository.AddressByUserId(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	if len(addressList) > 0 {
		user.Address = addressList[0]
	}

	return user, nil
}

func (s *userServ) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	var userRes model.User
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.UpdateUser(ctx, user)
		if errTx != nil {
			return errTx
		}
		userRes, errTx = s.User(ctx, user.ID)
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return model.User{}, err
	}

	return userRes, nil
}

func (s *userServ) GetUserRole(ctx context.Context, userID int64) (int, error) {
	user, err := s.userRepository.UserByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.Role, nil
}
