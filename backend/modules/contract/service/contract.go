package service

import (
	"context"
	"diploma/modules/contract/model"
	orderModel "diploma/modules/order/model"
	"fmt"
	"time"
)

type Repository interface {
	Create(ctx context.Context, contract *model.Contract) (int64, error)
	SignByParty(ctx context.Context, contractID int64, role int, signature string) error
	GetByID(ctx context.Context, id int64) (*model.Contract, error)
	MarkAsSigned(ctx context.Context, contractID int64) error
	GetByUser(ctx context.Context, userID int64) ([]*model.Contract, error) // 🔹 новый метод
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateContract(ctx context.Context, orderID, supplierID, customerID int64, content string) (int64, error) {
	contract := &model.Contract{
		OrderID:    orderID,
		SupplierID: supplierID,
		CustomerID: customerID,
		Content:    content,
		Status:     model.StatusCreated,
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, contract)
}

func (s *Service) SignContract(ctx context.Context, contractID int64, role int, signature string) error {
	contract, err := s.repo.GetByID(ctx, contractID)
	if err != nil {
		return err
	}

	switch role {
	case orderModel.CustomerRole:
		if !contract.SupplierSig.Valid {
			return fmt.Errorf("supplier has not signed the contract yet")
		}

	case orderModel.SupplierRole:
	default:
		return fmt.Errorf("unknown role")
	}

	err = s.repo.SignByParty(ctx, contractID, role, signature)
	if err != nil {
		return err
	}

	if role == orderModel.CustomerRole && contract.SupplierSig.Valid ||
		role == orderModel.SupplierRole && contract.CustomerSig.Valid {
		return s.repo.MarkAsSigned(ctx, contractID)
	}

	return nil
}

func (s *Service) GetContract(ctx context.Context, id int64) (*model.Contract, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetContractsByUser(ctx context.Context, userID int64) ([]*model.Contract, error) {
	return s.repo.GetByUser(ctx, userID)
}
