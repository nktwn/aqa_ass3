package service

import (
	"context"
	"diploma/modules/supplier/model"
)

func (s *SupplierService) SupplierListByIDList(ctx context.Context, id []int64) ([]model.Supplier, error) {
	return s.supplierRepo.SupplierListByIDList(ctx, id)
}

func (s *SupplierService) SupplierByID(ctx context.Context, id int64) (*model.Supplier, error) {
	supplier, err := s.supplierRepo.SupplierListByIDList(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	return &supplier[0], nil
}
