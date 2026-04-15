package converter

import (
	"diploma/modules/cart/model"
	modelRepo "diploma/modules/cart/repository/model"
)

func ToServiceCartFromRepo(input *modelRepo.Cart) *model.Cart {
	return &model.Cart{
		ID:         input.ID,
		CustomerID: input.CustomerID,
		Total:      input.Total,
	}
}

func ToServiceSupplierFromRepo(input []modelRepo.CartItem) []model.Supplier {
	res := make([]model.Supplier, 0, 1)
	supplierIndex := make(map[int64]int)
	for _, item := range input {
		index, ok := supplierIndex[item.SupplierID]
		if !ok {
			index = len(res)
			supplierIndex[item.SupplierID] = index
			res = append(res, model.Supplier{
				ID:          item.SupplierID,
				ProductList: make([]model.Product, 0, 1),
			})
		}
		res[index].ProductList = append(res[index].ProductList, model.Product{
			ID:       item.ProductID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Name:     item.ProductName,
			ImageUrl: item.ProductImageURL,
		})
	}
	return res
}
