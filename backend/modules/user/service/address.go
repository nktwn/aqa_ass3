package user

import (
	"context"

	"diploma/modules/user/model"
)

type IAddressRepo interface {
	CreateAddress(ctx context.Context, address model.Address) (int64, error)
	AddressByUserId(ctx context.Context, userID int64) ([]model.Address, error)
	GetById(ctx context.Context, id int64) (*model.Address, error)
	// Update(ctx context.Context, address *model.Address) error
	Delete(ctx context.Context, addressID int64) error
}

func (s *userServ) SetAddress(ctx context.Context, addr model.Address) (int64, error) {
	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateAddress(ctx, addr)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *userServ) Address(ctx context.Context, userID int64) ([]model.Address, error) {
	var addresList []model.Address

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		addresList, errTx = s.userRepository.AddressByUserId(ctx, userID)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return addresList, nil
}
