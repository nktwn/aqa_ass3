package model

type PutCartQuery struct {
	Quantity   int
	Price      int
	ProductID  int64
	SupplierID int64
	CustomerID int64
	CartID     int64
}

type Cart struct {
	ID         int64
	Total      int
	CustomerID int64
	Suppliers  []Supplier
}

type Supplier struct {
	OrderAmount        int
	TotalAmount        int
	FreeDeliveryAmount int
	DeliveryFee        int
	ID                 int64
	Name               string
	ProductList        []Product
}

type Product struct {
	Price    int
	Quantity int
	ID       int64
	Name     string
	ImageUrl string
}

type DeleteProductQuery struct {
	CustomerID int64
	ProductID  int64
	Quantity   int
	SupplierID int64
}


