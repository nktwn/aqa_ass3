package converter

import (
	"diploma/modules/supplier/model"
	repoModel "diploma/modules/supplier/repo/model"
)

func ToServiceSupplierFromRepo(supplier repoModel.Supplier) model.Supplier {
	return model.Supplier{
		ID:                 supplier.ID,
		Name:               supplier.Name,
		TotalAmount:        supplier.TotalAmount,
		OrderAmount:        supplier.OrderAmount,
		FreeDeliveryAmount: supplier.FreeDeliveryAmount,
		DeliveryFee:        supplier.DeliveryFee,
	}
}

func ToServiceSupplierListFromRepo(supplier []repoModel.Supplier) []model.Supplier {
	var suppliersModel []model.Supplier
	for _, s := range supplier {
		suppliersModel = append(suppliersModel, ToServiceSupplierFromRepo(s))
	}
	return suppliersModel

}
