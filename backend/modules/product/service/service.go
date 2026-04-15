package service

import (
	"context"
	"diploma/modules/product/model"
	"diploma/pkg/client/db"
	"diploma/pkg/service"
)

type ProductService struct {
	service.BaseService
	productRepository IProductRepository
	txManager         db.TxManager
}

func NewService(
	productRepository IProductRepository,
	txManager db.TxManager,
) *ProductService {
	return &ProductService{
		BaseService:       service.NewBaseService("product"),
		productRepository: productRepository,
		txManager:         txManager,
	}
}

type IProductRepository interface {
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error)
	GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error)
	GetProductList(ctx context.Context, query *model.ProductListQuery) ([]model.Product, error)
	GetTotalProducts(ctx context.Context, query *model.ProductListQuery) (int, error)
	GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error)
	GetProductSuggestions(ctx context.Context, query *model.SuggestionQuery) ([]string, error)
	GetReviewsByProduct(ctx context.Context, productID int64) ([]model.Review, error)
	UpsertReview(ctx context.Context, input *model.ReviewInput) error
	GetFavoriteProductIDsByUser(ctx context.Context, userID int64) ([]int64, error)
	IsFavorite(ctx context.Context, userID, productID int64) (bool, error)
	AddFavorite(ctx context.Context, userID, productID int64) error
	RemoveFavorite(ctx context.Context, userID, productID int64) error
	CreateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error)
	DeleteProduct(ctx context.Context, productID int64) error
}
