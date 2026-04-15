package supplier

import (
	"context"
	"diploma/modules/order/model"
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

func (a *SupplierClient) Supplier(ctx context.Context, id int64) (*model.Supplier, error) {
	supplier, err := a.supplierService.SupplierListByIDList(ctx, []int64{id})
	if err != nil {
		return nil, err
	}

	return &model.Supplier{
		ID:                 supplier[0].ID,
		Name:               supplier[0].Name,
		OrderAmount:        supplier[0].OrderAmount,
		DeliveryFee:        supplier[0].DeliveryFee,
		FreeDeliveryAmount: supplier[0].FreeDeliveryAmount,
	}, nil
}
