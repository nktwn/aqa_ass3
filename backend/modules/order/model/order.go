package model

import "time"

const (
	temp = iota
	Pending
	InProgress
	Completed
	Cancelled
)

type Order struct {
	ID              int64
	CustomerID      int64
	SupplierID      int64
	StatusID        int
	OrderDate       time.Time
	DeliveryAddress string
	DeliveryComment string
	PaymentID       string
	PaymentStatus   string
	ProductList     []*OrderProduct

	Supplier *Supplier
}

type OrderProduct struct {
	Quantity  int
	Price     int
	OrderID   int64
	ProductID int64
	Product   *Product
}

type Product struct {
	ID          int64
	Name        string
	Description string
	ImageUrl    string
}

type Supplier struct {
	ID                 int64
	Name               string
	OrderAmount        int
	FreeDeliveryAmount int
	DeliveryFee        int
}
