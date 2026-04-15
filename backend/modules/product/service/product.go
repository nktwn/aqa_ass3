package service

import (
	"context"
	"diploma/modules/product/model"
)

func (s *ProductService) Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error) {
	product, err := s.productRepository.GetProduct(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	productSupplierList, err := s.productRepository.GetSupplierProductListByProduct(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	reviews, err := s.productRepository.GetReviewsByProduct(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return &model.DetailedProduct{
		Product:             product,
		ProductSupplierList: productSupplierList,
		Reviews:             reviews,
	}, err

}

func (s *ProductService) ProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	return s.productRepository.GetProductPriceBySupplier(ctx, productID, supplierID)
}

func (s *ProductService) ProductInfo(ctx context.Context, id int64) (*model.Product, error) {
	return s.productRepository.GetProduct(ctx, id)
}

func (s *ProductService) FavoriteProducts(ctx context.Context, userID int64) ([]*model.Product, error) {
	ids, err := s.productRepository.GetFavoriteProductIDsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []*model.Product{}, nil
	}
	return s.productRepository.GetProductListByIDList(ctx, ids)
}

func (s *ProductService) FavoriteStatus(ctx context.Context, userID, productID int64) (bool, error) {
	return s.productRepository.IsFavorite(ctx, userID, productID)
}

func (s *ProductService) AddFavorite(ctx context.Context, userID, productID int64) error {
	return s.productRepository.AddFavorite(ctx, userID, productID)
}

func (s *ProductService) RemoveFavorite(ctx context.Context, userID, productID int64) error {
	return s.productRepository.RemoveFavorite(ctx, userID, productID)
}

func (s *ProductService) AddReview(ctx context.Context, input *model.ReviewInput) error {
	return s.productRepository.UpsertReview(ctx, input)
}

func (s *ProductService) CreateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error) {
	var product *model.Product
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		product, errTx = s.productRepository.CreateProduct(ctx, input)
		return errTx
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error) {
	var product *model.Product
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		product, errTx = s.productRepository.UpdateProduct(ctx, input)
		return errTx
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productID int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		return s.productRepository.DeleteProduct(ctx, productID)
	})
}
