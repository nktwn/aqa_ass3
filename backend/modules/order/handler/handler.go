package handler

import (
	"context"
	"diploma/modules/order/model"
)

type OrderHandler struct {
	service IOrderService
}

func NewHandler(service IOrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type IOrderService interface {
	ICreateOrderService
	GetOrderByID(ctx context.Context, userID int64, role int, orderID int64) (*model.Order, error)
	UpdateOrderStatusBySupplier(ctx context.Context, supplierID int64, orderID int64, newStatusID int) error
	CancelOrderByCustomer(ctx context.Context, customerID int64, orderID int64) error
}
