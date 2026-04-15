package model

type ProductInput struct {
	ID int64 `json:"id"`
}

type ProductResponse struct {
	DetailedProduct *DetailedProduct `json:"product"`
}

type Product struct {
	ID                    int64           `json:"id"`
	Name                  string          `json:"name"`
	ImageUrl              string          `json:"image"`
	RatingAverage         float64         `json:"rating_average"`
	RatingCount           int             `json:"rating_count"`
	LowestProductSupplier ProductSupplier `json:"lowest_product_supplier"`
}

type ProductSupplier struct {
	Price      int      `json:"price"`
	SellAmount int      `json:"sell_amount"`
	Supplier   Supplier `json:"supplier"`
}

type Supplier struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	OrderAmount        int    `json:"order_amount"`
	FreeDeliveryAmount int    `json:"free_delivery_amount"`
	DeliveryFee        int    `json:"delivery_fee"`
}

type DetailedProduct struct {
	*Product            `json:"product"`
	ProductSupplierList []ProductSupplier `json:"suppliers"`
	Reviews             []Review          `json:"reviews"`
}

type ErrorResponse struct {
	Err string `json:"error"`
}

type Review struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type ReviewInput struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type FavoriteStatusResponse struct {
	IsFavorite bool `json:"is_favorite"`
}

type FavoriteListResponse struct {
	ProductList []Product `json:"product_list"`
}

type AdminProductInput struct {
	Name     string `json:"name" binding:"required"`
	ImageUrl string `json:"image"`
	GTIN     int64  `json:"gtin"`
}

type ProductMessageResponse struct {
	Message string `json:"message"`
}

type AdminProductResponse struct {
	Product *Product `json:"product"`
}
