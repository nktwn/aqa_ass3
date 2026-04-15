package converter

import (
	"diploma/modules/product/model"
	repoModel "diploma/modules/product/repository/product/model"
)

func toProductFromRepo(p repoModel.Product) model.Product {
	return model.Product{
		ID:                    p.ID,
		Name:                  p.Name,
		ImageUrl:              p.ImageUrl,
		GTIN:                  p.GTIN,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
		RatingAverage:         p.RatingAverage,
		RatingCount:           p.RatingCount,
		LowestProductSupplier: ToProductSupplierFromRepo(p.LowestSupplier),
	}
}

func ToProductListFromRepo(list []repoModel.Product) []model.Product {
	result := make([]model.Product, len(list))
	for i, p := range list {
		result[i] = toProductFromRepo(p)
	}
	return result
}

// ToProductFromRepo converts a repository-level Product to a domain-level Product.
func ToProductFromRepo(p repoModel.Product) *model.Product {
	return &model.Product{
		ID:                    p.ID,
		Name:                  p.Name,
		ImageUrl:              p.ImageUrl,
		GTIN:                  p.GTIN,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
		RatingAverage:         p.RatingAverage,
		RatingCount:           p.RatingCount,
		LowestProductSupplier: ToProductSupplierFromRepo(p.LowestSupplier),
	}
}

func ToProductSupplierFromRepo(sp repoModel.ProductSupplier) model.ProductSupplier {
	return model.ProductSupplier{
		Price:      sp.Price,
		SellAmount: sp.SellAmount,
		Supplier:   model.Supplier(sp.Supplier),
	}
}

func ToSupplierFromRepo(s repoModel.Supplier) model.Supplier {
	return model.Supplier{
		ID:                 s.ID,
		Name:               s.Name,
		OrderAmount:        s.OrderAmount,
		FreeDeliveryAmount: s.FreeDeliveryAmount,
		DeliveryFee:        s.DeliveryFee,
	}
}
