package model

import "time"

type Product struct {
	ID             int64
	GTIN           int64
	Name           string
	ImageUrl       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RatingAverage  float64
	RatingCount    int
	LowestSupplier ProductSupplier
}

type ProductSupplier struct {
	Price      int
	SellAmount int
	Supplier   Supplier
}

type Supplier struct {
	ID                 int64
	Name               string
	OrderAmount        int
	FreeDeliveryAmount int
	DeliveryFee        int
}

type ProductListQuery struct {
	Offset    int
	Limit     int
	Search    string
	MinPrice  int
	MaxPrice  int
	SortBy    string
	SortOrder string
}
