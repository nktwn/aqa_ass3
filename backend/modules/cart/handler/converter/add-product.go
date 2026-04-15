package converter

import (
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
)

func ToServiceCardInputFromAPI(input *modelApi.AddProductToCartInput) *model.PutCartQuery {
	return &model.PutCartQuery{
		CustomerID: input.CustomerID,
		ProductID:  input.ProductID,
		Quantity:   input.Quantity,
		SupplierID: input.SupplierID,
	}
}

func ToAPIGetCartFromService(card *model.Cart) *modelApi.GetCartResponse {
	return &modelApi.GetCartResponse{
		Total:      card.Total,
		CustomerID: card.CustomerID,
		Suppliers:  ToAPISuppliersFromService(card.Suppliers),
	}
}

func ToAPISuppliersFromService(suppliers []model.Supplier) []modelApi.Supplier {
	apiSuppliers := make([]modelApi.Supplier, len(suppliers))
	for i, supplier := range suppliers {
		apiSuppliers[i] = modelApi.Supplier{
			ID:                 supplier.ID,
			Name:               supplier.Name,
			OrderAmount:        supplier.OrderAmount,
			TotalAmount:        supplier.TotalAmount,
			FreeDeliveryAmount: supplier.FreeDeliveryAmount,
			DeliveryFee:        supplier.DeliveryFee,
			ProductList:        ToAPIProductsFromService(supplier.ProductList),
		}
	}
	return apiSuppliers
}

func ToAPIProductsFromService(products []model.Product) []modelApi.Product {
	apiProducts := make([]modelApi.Product, len(products))
	for i, product := range products {
		apiProducts[i] = modelApi.Product{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quantity,
			ImageUrl:    product.ImageUrl,
		}
	}
	return apiProducts
}
