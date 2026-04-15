package model

type ProductListInput struct {
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Search    string `json:"search"`
	MinPrice  int    `json:"min_price"`
	MaxPrice  int    `json:"max_price"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

type ProductListResponse struct {
	ProductList []Product `json:"product_list"`
	Total       int       `json:"total"`
}

type ProductSuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
}
