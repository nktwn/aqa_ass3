package model

type Supplier struct {
	OrderAmount        int
	TotalAmount        int
	FreeDeliveryAmount int
	DeliveryFee        int
	ID                 int64
	Name               string
}
