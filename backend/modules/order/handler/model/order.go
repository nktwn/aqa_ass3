package model

import "errors"

var (
	ErrUnauthorized = errors.New("api: unauthorized")
)

type ErrorResponse struct {
	Err string `json:"error"`
}

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID              int64      `json:"id"`
	Status          string     `json:"status"`
	OrderDate       string     `json:"order_date"`
	DeliveryAddress string     `json:"delivery_address"`
	DeliveryComment string     `json:"delivery_comment"`
	PaymentStatus   string     `json:"payment_status"`
	Supplier        *Supplier  `json:"supplier"`
	ProductList     []*Product `json:"product_list"`
}
type Product struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}
type Supplier struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateOrderStatusRequest struct {
	OrderID     int64 `json:"order_id" binding:"required"`
	NewStatusID int   `json:"new_status_id" binding:"required"`
}

type CancelOrderRequest struct {
	OrderID int64 `json:"order_id" binding:"required"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
