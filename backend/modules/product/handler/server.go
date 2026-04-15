package handler

import (
	"context"
	"diploma/modules/product/model"
)

type CatalogHandler struct {
	service IProductService
}

func NewHandler(service IProductService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

type IProductService interface {
	ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error)
	Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error)
	FeaturedProducts(ctx context.Context, limit int) (*model.ProductList, error)
	Suggestions(ctx context.Context, query *model.SuggestionQuery) ([]string, error)
	FavoriteProducts(ctx context.Context, userID int64) ([]*model.Product, error)
	FavoriteStatus(ctx context.Context, userID, productID int64) (bool, error)
	AddFavorite(ctx context.Context, userID, productID int64) error
	RemoveFavorite(ctx context.Context, userID, productID int64) error
	AddReview(ctx context.Context, input *model.ReviewInput) error
	CreateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error)
	DeleteProduct(ctx context.Context, productID int64) error
}
