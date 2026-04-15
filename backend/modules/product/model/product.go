package model

import "time"

type ProductListQuery struct {
	Offset    int
	Limit     int
	Search    string
	MinPrice  int
	MaxPrice  int
	SortBy    string
	SortOrder string
}

type ProductList struct {
	Products []Product
	Total    int
}

type ProductQuery struct {
	ID int64
}

type DetailedProduct struct {
	*Product
	ProductSupplierList []ProductSupplier
	Reviews             []Review
}

type Product struct {
	ID             int64
	GTIN           int64
	Name           string
	ImageUrl       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RatingAverage  float64
	RatingCount    int
	LowestProductSupplier ProductSupplier
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

type SuggestionQuery struct {
	Query string
	Limit int
}

type Review struct {
	ID        int64
	UserID    int64
	UserName  string
	Rating    int
	Comment   string
	CreatedAt time.Time
}

type ReviewInput struct {
	ProductID int64
	UserID    int64
	Rating    int
	Comment   string
}

type AdminProductInput struct {
	ID       int64
	Name     string
	ImageUrl string
	GTIN     int64
}
