package model

import "errors"

var (
	ErrUnauthorized = errors.New("api: unauthorized")
)

type AddProductToCartInput struct {
	Quantity   int   `json:"quantity"`
	ProductID  int64 `json:"product_id"`
	SupplierID int64 `json:"supplier_id"`
	CustomerID int64 `json:"-"`
}

type AddProductToCardResponse struct {
	Status string `json:"status"`
}

type GetCardInput struct {
	CustomerID int64 `json:"CustomerID"`
}

type GetCartResponse struct {
	Total      int        `json:"total"`
	CustomerID int64      `json:"customer_id"`
	Suppliers  []Supplier `json:"suppliers"`
}

type Supplier struct {
	OrderAmount        int       `json:"order_amount"`
	TotalAmount        int       `json:"total_amount"`
	FreeDeliveryAmount int       `json:"free_delivery_amount"`
	DeliveryFee        int       `json:"delivery_fee"`
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	ProductList        []Product `json:"product_list"`
}

type Product struct {
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image"`
}

type DeleteProductFromCartInput struct {
	CustomerID int64 `json:"customer_id"`
	ProductID  int64 `json:"product_id"`
	SupplierID int64 `json:"supplier_id"`
	Quantity   int   `json:"quantity"`
}

type ErrorResponse struct {
	Err string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
