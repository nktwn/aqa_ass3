package service

import (
	"context"
	"diploma/modules/supplier/model"
	"diploma/pkg/client/db"
)

type SupplierService struct {
	supplierRepo ISupplierRepository
	txManager    db.TxManager
}

func NewService(repo ISupplierRepository, tx db.TxManager) *SupplierService {
	return &SupplierService{
			supplierRepo: repo,
		txManager:    tx,
	}

}

type ISupplierRepository interface {
	SupplierListByIDList(ctx context.Context, id []int64) ([]model.Supplier, error)
}
