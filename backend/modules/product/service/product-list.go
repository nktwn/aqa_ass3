package service

import (
	"context"
	"diploma/modules/product/model"
	"strings"

	"go.uber.org/zap"
)

func (s *ProductService) ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error) {
	s.LogInfo(ctx, "Fetching product list",
		zap.Int("offset", query.Offset),
		zap.Int("limit", query.Limit),
	)

	productList, err := s.productRepository.GetProductList(ctx, query)
	if err != nil {
		s.LogError(ctx, "Failed to get product list", err)
		return nil, err
	}

	total, err := s.productRepository.GetTotalProducts(ctx, query)
	if err != nil {
		s.LogError(ctx, "Failed to get total products count", err)
		return nil, err
	}

	result := &model.ProductList{
		Products: productList,
		Total:    total,
	}

	s.LogInfo(ctx, "Successfully fetched product list",
		zap.Int("total_products", total),
		zap.Int("returned_products", len(productList)),
	)

	return result, nil
}

func (s *ProductService) FeaturedProducts(ctx context.Context, limit int) (*model.ProductList, error) {
	if limit <= 0 {
		limit = 8
	}

	return s.ProductList(ctx, &model.ProductListQuery{
		Limit:     limit,
		Offset:    0,
		SortBy:    "newest",
		SortOrder: "desc",
	})
}

func (s *ProductService) Suggestions(ctx context.Context, query *model.SuggestionQuery) ([]string, error) {
	if query == nil {
		query = &model.SuggestionQuery{}
	}

	query.Query = strings.TrimSpace(query.Query)
	if query.Query == "" {
		return []string{}, nil
	}

	if query.Limit <= 0 {
		query.Limit = 5
	}

	suggestions, err := s.productRepository.GetProductSuggestions(ctx, query)
	if err != nil {
		s.LogError(ctx, "Failed to get product suggestions", err)
		return nil, err
	}

	return suggestions, nil
}

func (s *ProductService) ProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
	s.LogInfo(ctx, "Fetching products by ID list",
		zap.Int("product_count", len(idList)),
	)

	productList, err := s.productRepository.GetProductListByIDList(ctx, idList)
	if err != nil {
		s.LogError(ctx, "Failed to get products by ID list", err,
			zap.Any("product_ids", idList),
		)
		return nil, err
	}

	s.LogInfo(ctx, "Successfully fetched products by ID list",
		zap.Int("requested_count", len(idList)),
		zap.Int("found_count", len(productList)),
	)

	return productList, nil
}
