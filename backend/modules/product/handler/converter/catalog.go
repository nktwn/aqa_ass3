package converter

import (
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
)

func ToServieProductQueryFromApi(input modelApi.ProductInput) *model.ProductQuery {
	return &model.ProductQuery{
		ID: input.ID,
	}
}

func ToApiProductResponeFromService(product *model.DetailedProduct) *modelApi.ProductResponse {
	return &modelApi.ProductResponse{
		DetailedProduct: ToApiDetailedProductFromSerivce(product),
	}
}

func ToApiDetailedProductFromSerivce(dProudct *model.DetailedProduct) *modelApi.DetailedProduct {
	return &modelApi.DetailedProduct{
		Product:             ToAPIProductFromService(dProudct.Product),
		ProductSupplierList: ToApiProductSupplierListFromService(dProudct.ProductSupplierList),
		Reviews:             ToAPIReviewsFromService(dProudct.Reviews),
	}
}

func ToApiProductSupplierListFromService(productSupplierList []model.ProductSupplier) []modelApi.ProductSupplier {
	res := make([]modelApi.ProductSupplier, 0, len(productSupplierList))

	for _, ps := range productSupplierList {
		res = append(res, ToAPIProductSupplierFromService(ps))
	}

	return res
}

func ToAPIProductFromService(product *model.Product) *modelApi.Product {
	return &modelApi.Product{
		ID:                    product.ID,
		Name:                  product.Name,
		ImageUrl:              product.ImageUrl,
		RatingAverage:         product.RatingAverage,
		RatingCount:           product.RatingCount,
		LowestProductSupplier: ToAPIProductSupplierFromService(product.LowestProductSupplier),
	}
}

func ToAPIReviewsFromService(reviews []model.Review) []modelApi.Review {
	result := make([]modelApi.Review, 0, len(reviews))
	for _, review := range reviews {
		result = append(result, modelApi.Review{
			ID:        review.ID,
			UserID:    review.UserID,
			UserName:  review.UserName,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result
}

// ConvertServiceToAPISuppSlierInfo преобразует информацию о поставщике из сервиса в API.
func ToAPIProductSupplierFromService(ps model.ProductSupplier) modelApi.ProductSupplier {
	return modelApi.ProductSupplier{
		Price:      ps.Price,
		SellAmount: ps.SellAmount,
		Supplier:   ToAPISupplierFromService(ps.Supplier),
	}
}

func ToAPISupplierFromService(supplier model.Supplier) modelApi.Supplier {
	return modelApi.Supplier{
		ID:                 supplier.ID,
		Name:               supplier.Name,
		OrderAmount:        supplier.OrderAmount,
		FreeDeliveryAmount: supplier.FreeDeliveryAmount,
		DeliveryFee:        supplier.DeliveryFee,
	}
}

// ConvertAPIToServiceSupplierInfo преобразует информацию о поставщике из API в сервис.
// func ToServiceSupplierInfoFromAPI(apiSupplierInfo modelApi.ProductSupplierInfo) model.ProductSupplierInfo {
// 	return model.ProductSupplierInfo{
// 		SupplierID:                apiSupplierInfo.SupplierID,
// 		Name:                      apiSupplierInfo.Name,
// 		MinimumFreeDeliveryAmount: apiSupplierInfo.MinimumFreeDeliveryAmount,
// 		DeliveryFee:               apiSupplierInfo.DeliveryFee,
// 	}
// }
