package supplier

import (
	"context"
	"diploma/modules/cart/model"
	supplierModel "diploma/modules/supplier/model"
)

type SupplierClient struct {
	supplierService ISupplierService
}

func NewClient(supplierService ISupplierService) *SupplierClient {
	return &SupplierClient{supplierService: supplierService}
}

type ISupplierService interface {
	SupplierListByIDList(ctx context.Context, idList []int64) ([]supplierModel.Supplier, error)
}

func (a *SupplierClient) SupplierListByIDList(ctx context.Context, IDList []int64) ([]model.Supplier, error) {

	suppliers, err := a.supplierService.SupplierListByIDList(ctx, IDList)
	if err != nil {
		return nil, err
	}
	var res []model.Supplier
	for _, supplier := range suppliers {
		res = append(res, model.Supplier{
			ID:                 supplier.ID,
			Name:               supplier.Name,
			OrderAmount:        supplier.OrderAmount,
			DeliveryFee:        supplier.DeliveryFee,
			FreeDeliveryAmount: supplier.FreeDeliveryAmount,
		})
	}
	return res, nil
}
