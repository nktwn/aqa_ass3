package model

type Cart struct {
	ID         int64
	CustomerID int64
	Total      int
	Items      []CartItem
}

type CartItem struct {
	ID              int64
	ProductID       int64
	ProductName     string
	ProductImageURL string
	SupplierID      int64
	Quantity        int
	Price           int
}
