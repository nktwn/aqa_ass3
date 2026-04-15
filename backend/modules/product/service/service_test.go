package service

import (
	"context"
	"diploma/internal/testutils"
	"diploma/modules/product/model"
	"diploma/pkg/client/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Product), args.Error(1)
}
func (m *mockProductRepository) GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]model.ProductSupplier), args.Error(1)
}
func (m *mockProductRepository) GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
	args := m.Called(ctx, idList)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*model.Product), args.Error(1)
}
func (m *mockProductRepository) GetProductList(ctx context.Context, query *model.ProductListQuery) ([]model.Product, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]model.Product), args.Error(1)
}
func (m *mockProductRepository) GetTotalProducts(ctx context.Context, query *model.ProductListQuery) (int, error) {
	args := m.Called(ctx, query)
	return args.Int(0), args.Error(1)
}
func (m *mockProductRepository) GetProductSuggestions(ctx context.Context, query *model.SuggestionQuery) ([]string, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]string), args.Error(1)
}
func (m *mockProductRepository) GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	args := m.Called(ctx, productID, supplierID)
	return args.Int(0), args.Error(1)
}
func (m *mockProductRepository) GetReviewsByProduct(ctx context.Context, productID int64) ([]model.Review, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]model.Review), args.Error(1)
}
func (m *mockProductRepository) UpsertReview(ctx context.Context, input *model.ReviewInput) error {
	return m.Called(ctx, input).Error(0)
}
func (m *mockProductRepository) GetFavoriteProductIDsByUser(ctx context.Context, userID int64) ([]int64, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]int64), args.Error(1)
}
func (m *mockProductRepository) IsFavorite(ctx context.Context, userID, productID int64) (bool, error) {
	args := m.Called(ctx, userID, productID)
	return args.Bool(0), args.Error(1)
}
func (m *mockProductRepository) AddFavorite(ctx context.Context, userID, productID int64) error {
	return m.Called(ctx, userID, productID).Error(0)
}
func (m *mockProductRepository) RemoveFavorite(ctx context.Context, userID, productID int64) error {
	return m.Called(ctx, userID, productID).Error(0)
}
func (m *mockProductRepository) CreateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Product), args.Error(1)
}
func (m *mockProductRepository) UpdateProduct(ctx context.Context, input *model.AdminProductInput) (*model.Product, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Product), args.Error(1)
}
func (m *mockProductRepository) DeleteProduct(ctx context.Context, productID int64) error {
	return m.Called(ctx, productID).Error(0)
}

type mockTxManager struct { mock.Mock }
func (m *mockTxManager) ReadCommitted(ctx context.Context, h db.Handler) error {
	args := m.Called(ctx, h)
	if h != nil {
		if err := h(ctx); err != nil { return err }
	}
	return args.Error(0)
}

type ProductServiceTestSuite struct {
	suite.Suite
	service    *ProductService
	repository *mockProductRepository
	txManager  *mockTxManager
	helper     *testutils.AssertTestHelper
}

func TestProductService(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupTest() {
	s.repository = new(mockProductRepository)
	s.txManager = new(mockTxManager)
	s.helper = testutils.NewAssertTestHelper(s.T())
	s.service = NewService(s.repository, s.txManager)
}

func (s *ProductServiceTestSuite) TestProductList_Success() {
	query := &model.ProductListQuery{Offset: 0, Limit: 10}
	now := time.Now()
	expectedProducts := []model.Product{{ID: 1, Name: "Product 1", GTIN: 123456789, ImageUrl: "http://example.com/image1.jpg", CreatedAt: now, UpdatedAt: now}}

	s.repository.On("GetProductList", mock.Anything, query).Return(expectedProducts, nil).Once()
	s.repository.On("GetTotalProducts", mock.Anything, query).Return(1, nil).Once()

	result, err := s.service.ProductList(context.Background(), query)

	s.helper.AssertNoError(err)
	assert.Equal(s.T(), 1, result.Total)
	assert.Len(s.T(), result.Products, 1)
	assert.Equal(s.T(), int64(1), result.Products[0].ID)
}

func (s *ProductServiceTestSuite) TestSuggestions_TrimsQueryAndDefaultsLimit() {
	query := &model.SuggestionQuery{Query: "  coffee  "}
	s.repository.On("GetProductSuggestions", mock.Anything, mock.MatchedBy(func(input *model.SuggestionQuery) bool {
		return input.Query == "coffee" && input.Limit == 5
	})).Return([]string{"Arabica Coffee Beans 1kg"}, nil).Once()

	result, err := s.service.Suggestions(context.Background(), query)

	s.helper.AssertNoError(err)
	assert.Equal(s.T(), []string{"Arabica Coffee Beans 1kg"}, result)
}

func (s *ProductServiceTestSuite) TestFavoriteProducts_EmptyList() {
	s.repository.On("GetFavoriteProductIDsByUser", mock.Anything, int64(1)).Return([]int64{}, nil).Once()

	products, err := s.service.FavoriteProducts(context.Background(), 1)

	s.helper.AssertNoError(err)
	assert.Empty(s.T(), products)
	s.repository.AssertNotCalled(s.T(), "GetProductListByIDList")
}

func (s *ProductServiceTestSuite) TestCreateProduct_UsesTransaction() {
	input := &model.AdminProductInput{Name: "New product", GTIN: 12345}
	expected := &model.Product{ID: 10, Name: "New product", GTIN: 12345}

	s.txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	s.repository.On("CreateProduct", mock.Anything, input).Return(expected, nil).Once()

	product, err := s.service.CreateProduct(context.Background(), input)

	s.helper.AssertNoError(err)
	assert.Equal(s.T(), int64(10), product.ID)
	assert.Equal(s.T(), "New product", product.Name)
}

func TestSuggestions_EmptyQueryReturnsEmpty(t *testing.T) {
	t.Parallel()

	service := NewService(new(mockProductRepository), new(mockTxManager))
	result, err := service.Suggestions(context.Background(), &model.SuggestionQuery{Query: "   "})
	require.NoError(t, err)
	assert.Empty(t, result)
}
