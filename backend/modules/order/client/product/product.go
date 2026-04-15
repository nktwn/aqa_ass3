package product

import (
	"context"
	"diploma/modules/order/model"
	clientModel "diploma/modules/product/model"
)

type ProductClient struct {
	supplierService IProductSerivce
}

func NewClient(supplierService IProductSerivce) *ProductClient {
	return &ProductClient{supplierService: supplierService}
}

type IProductSerivce interface {
	ProductInfo(ctx context.Context, id int64) (*clientModel.Product, error)
}

func (a *ProductClient) Product(ctx context.Context, id int64) (*model.Product, error) {
	product, err := a.supplierService.ProductInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	// Convert clientModel.Product to model.Product
	return &model.Product{
		ID:       product.ID,
		Name:     product.Name,
		ImageUrl: product.ImageUrl,
	}, nil

}
