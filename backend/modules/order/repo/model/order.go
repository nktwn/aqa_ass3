package model

import "time"

type Order struct {
	ID              int64
	SupplierID      int64
	CustomerID      int64
	StatusID        int
	OrderDate       time.Time
	DeliveryAddress string
	DeliveryComment string
	PaymentID       string
	PaymentStatus   string
}
