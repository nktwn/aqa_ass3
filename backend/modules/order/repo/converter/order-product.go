package converter

import (
	"diploma/modules/order/model"
	modelRepo "diploma/modules/order/repo/model"
)

func ToServiceOrderProductFromRepo(op modelRepo.OrderProduct) *model.OrderProduct {
	return &model.OrderProduct{
		Quantity:  op.Quantity,
		Price:     op.Price,
		OrderID:   op.OrderID,
		ProductID: op.ProductID,
	}
}

func ToServiceOrderFromRepo(o *modelRepo.Order) *model.Order {
	return &model.Order{
		ID:              o.ID,
		CustomerID:      o.CustomerID,
		SupplierID:      o.SupplierID,
		StatusID:        o.StatusID,
		OrderDate:       o.OrderDate,
		DeliveryAddress: o.DeliveryAddress,
		DeliveryComment: o.DeliveryComment,
		PaymentID:       o.PaymentID,
		PaymentStatus:   o.PaymentStatus,
	}
}
