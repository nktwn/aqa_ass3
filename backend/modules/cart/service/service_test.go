package service

import (
	"context"
	"diploma/internal/testutils"
	"diploma/modules/cart/model"
	userModel "diploma/modules/user/model"
	"diploma/pkg/client/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockCartRepository struct {
	mock.Mock
}

func (m *mockCartRepository) UpdateItemQuantity(ctx context.Context, cartId, productId, supplierId int64, quantity int) error {
	args := m.Called(ctx, cartId, productId, supplierId, quantity)
	return args.Error(0)
}

func (m *mockCartRepository) ItemQuantity(ctx context.Context, cartId, productId, supplierId int64) (int, error) {
	args := m.Called(ctx, cartId, productId, supplierId)
	return args.Int(0), args.Error(1)
}

func (m *mockCartRepository) Cart(ctx context.Context, userID int64) (*model.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *mockCartRepository) CreateCart(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockCartRepository) AddItem(ctx context.Context, input *model.PutCartQuery) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *mockCartRepository) UpdateCartTotal(ctx context.Context, cartID int64, total int) error {
	args := m.Called(ctx, cartID, total)
	return args.Error(0)
}

func (m *mockCartRepository) DeleteCart(ctx context.Context, cartID int64) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

func (m *mockCartRepository) GetCartItems(ctx context.Context, cartID int64) ([]model.Supplier, error) {
	args := m.Called(ctx, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Supplier), args.Error(1)
}

func (m *mockCartRepository) DeleteCartItems(ctx context.Context, cartID int64) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

func (m *mockCartRepository) DeleteItem(ctx context.Context, cartID, productId, supplierId int64) error {
	args := m.Called(ctx, cartID, productId, supplierId)
	return args.Error(0)
}

type mockProductService struct { mock.Mock }
func (m *mockProductService) ProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	args := m.Called(ctx, productID, supplierID)
	return args.Int(0), args.Error(1)
}

type mockSupplierClient struct { mock.Mock }
func (m *mockSupplierClient) SupplierListByIDList(ctx context.Context, idList []int64) ([]model.Supplier, error) {
	args := m.Called(ctx, idList)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]model.Supplier), args.Error(1)
}

type mockOrderClient struct { mock.Mock }
func (m *mockOrderClient) CreateOrder(ctx context.Context, cart *model.Cart, address model.DeliveryAddress, paymentID string) error {
	args := m.Called(ctx, cart, address, paymentID)
	return args.Error(0)
}

type mockAddressBook struct { mock.Mock }
func (m *mockAddressBook) GetById(ctx context.Context, id int64) (*userModel.Address, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*userModel.Address), args.Error(1)
}

type mockPaymentClient struct { mock.Mock }
func (m *mockPaymentClient) PaymentRequest(orderID, amount, currency, description string) (model.CheckoutResponse, error) {
	args := m.Called(orderID, amount, currency, description)
	return args.Get(0).(model.CheckoutResponse), args.Error(1)
}

type mockRedis struct { mock.Mock }
func (m *mockRedis) SavePaymentOrder(ctx context.Context, paymentOrder model.PaymentOrder) error {
	args := m.Called(ctx, paymentOrder)
	return args.Error(0)
}
func (m *mockRedis) PaymentOrder(ctx context.Context, orderID string) (model.PaymentOrder, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil { return model.PaymentOrder{}, args.Error(1) }
	return args.Get(0).(model.PaymentOrder), args.Error(1)
}

type mockTxManager struct { mock.Mock }
func (m *mockTxManager) ReadCommitted(ctx context.Context, h db.Handler) error {
	args := m.Called(ctx, h)
	if h != nil {
		if err := h(ctx); err != nil {
			return err
		}
	}
	return args.Error(0)
}

type CartServiceTestSuite struct {
	suite.Suite
	service         *cartServ
	cartRepo        *mockCartRepository
	productService  *mockProductService
	supplierService *mockSupplierClient
	orderService    *mockOrderClient
	addressBook     *mockAddressBook
	paymentClient   *mockPaymentClient
	redis           *mockRedis
	txManager       *mockTxManager
	helper          *testutils.AssertTestHelper
}

func TestCartService(t *testing.T) {
	suite.Run(t, new(CartServiceTestSuite))
}

func (s *CartServiceTestSuite) SetupTest() {
	s.cartRepo = new(mockCartRepository)
	s.productService = new(mockProductService)
	s.supplierService = new(mockSupplierClient)
	s.orderService = new(mockOrderClient)
	s.addressBook = new(mockAddressBook)
	s.paymentClient = new(mockPaymentClient)
	s.redis = new(mockRedis)
	s.txManager = new(mockTxManager)
	s.helper = testutils.NewAssertTestHelper(s.T())

	s.service = NewService(
		s.cartRepo,
		s.productService,
		s.supplierService,
		s.orderService,
		s.addressBook,
		s.paymentClient,
		s.redis,
		s.txManager,
	)
}

func (s *CartServiceTestSuite) TestAddProductToCard_CreatesCartAndAddsItem() {
	query := &model.PutCartQuery{CustomerID: 1, ProductID: 101, SupplierID: 201, Quantity: 2}

	s.txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	s.cartRepo.On("Cart", mock.Anything, int64(1)).Return(nil, model.ErrNoRows).Once()
	s.cartRepo.On("CreateCart", mock.Anything, int64(1)).Return(int64(77), nil).Once()
	s.productService.On("ProductPriceBySupplier", mock.Anything, int64(101), int64(201)).Return(1200, nil).Once()
	s.cartRepo.On("ItemQuantity", mock.Anything, int64(77), int64(101), int64(201)).Return(0, model.ErrNoRows).Once()
	s.cartRepo.On("AddItem", mock.Anything, mock.MatchedBy(func(input *model.PutCartQuery) bool {
		return input.CartID == 77 && input.Price == 1200 && input.Quantity == 2
	})).Return(nil).Once()
	s.cartRepo.On("UpdateCartTotal", mock.Anything, int64(77), 2400).Return(nil).Once()

	err := s.service.AddProductToCard(context.Background(), query)

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(int64(77), query.CartID)
	s.helper.AssertEqual(1200, query.Price)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *CartServiceTestSuite) TestCheckout_SavesPaymentOrder() {
	ctx := context.Background()
	address := &userModel.Address{ID: 10, UserID: 1, Street: "Almaty, Abay 10", Description: "Home"}
	cartRecord := &model.Cart{ID: 7, CustomerID: 1}
	cartItems := []model.Supplier{{
		ID:                 201,
		OrderAmount:        2500,
		FreeDeliveryAmount: 10000,
		DeliveryFee:        500,
		ProductList: []model.Product{{ID: 101, Price: 2500, Quantity: 1}},
	}}
	supplierDetails := []model.Supplier{{
		ID:                 201,
		Name:               "Roastery One",
		OrderAmount:        2500,
		FreeDeliveryAmount: 10000,
		DeliveryFee:        500,
	}}

	s.txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Twice()
	s.addressBook.On("GetById", mock.Anything, int64(10)).Return(address, nil).Once()
	s.cartRepo.On("Cart", mock.Anything, int64(1)).Return(cartRecord, nil).Once()
	s.cartRepo.On("GetCartItems", mock.Anything, int64(7)).Return(cartItems, nil).Once()
	s.supplierService.On("SupplierListByIDList", mock.Anything, []int64{201}).Return(supplierDetails, nil).Once()
	s.paymentClient.On("PaymentRequest", mock.AnythingOfType("string"), "300000", "USD", "Payment for order").
		Return(model.CheckoutResponse{CheckoutURL: "https://pay", PaymentID: "payment-1", ResponseStatus: "ok"}, nil).Once()
	s.redis.On("SavePaymentOrder", mock.Anything, mock.MatchedBy(func(order model.PaymentOrder) bool {
		return order.Cart.ID == 7 && order.DeliveryAddress.ID == 10 && order.PaymentID == "payment-1"
	})).Return(nil).Once()

	checkout, err := s.service.Checkout(ctx, 1, model.CheckoutInput{AddressID: 10})

	s.helper.AssertNoError(err)
	assert.Equal(s.T(), "https://pay", checkout.CheckoutURL)
	assert.Equal(s.T(), "payment-1", checkout.PaymentID)
	assert.NotEmpty(s.T(), checkout.OrderID)
}

func (s *CartServiceTestSuite) TestCheckout_RejectsForeignAddress() {
	s.txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	s.addressBook.On("GetById", mock.Anything, int64(10)).Return(&userModel.Address{ID: 10, UserID: 99}, nil).Once()

	_, err := s.service.Checkout(context.Background(), 1, model.CheckoutInput{AddressID: 10})

	s.helper.AssertError(err)
	assert.Equal(s.T(), "address does not belong to current user", err.Error())
}

func (s *CartServiceTestSuite) TestCommitCheckout_Approved_CreatesOrderAndClearsCart() {
	paymentOrder := model.PaymentOrder{
		ID: "1-5-order",
		Cart: model.Cart{ID: 5, CustomerID: 1, Suppliers: []model.Supplier{{ID: 201}}},
		DeliveryAddress: model.DeliveryAddress{ID: 11, Street: "Almaty, Satpayev 5", Description: "Office"},
		CheckoutResponse: model.CheckoutResponse{PaymentID: "payment-55"},
	}

	s.redis.On("PaymentOrder", mock.Anything, "1-5-order").Return(paymentOrder, nil).Once()
	s.txManager.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	s.orderService.On("CreateOrder", mock.Anything, &paymentOrder.Cart, paymentOrder.DeliveryAddress, "payment-55").Return(nil).Once()
	s.cartRepo.On("DeleteCartItems", mock.Anything, int64(5)).Return(nil).Once()
	s.cartRepo.On("DeleteCart", mock.Anything, int64(5)).Return(nil).Once()

	err := s.service.CommitCheckout(context.Background(), model.CommitCheckout{OrderID: "1-5-order", PaymentStatus: model.PaymentStatusApproved})

	s.helper.AssertNoError(err)
	s.orderService.AssertExpectations(s.T())
	s.cartRepo.AssertExpectations(s.T())
}

func TestCartBusinessRules(t *testing.T) {
	t.Parallel()

	t.Run("getTotalSupplier adds delivery fee below threshold", func(t *testing.T) {
		t.Parallel()
		total := getTotalSupplier([]model.Product{{Price: 1000, Quantity: 2}}, model.Supplier{FreeDeliveryAmount: 5000, DeliveryFee: 700})
		assert.Equal(t, 2700, total)
	})

	t.Run("checkCartForCheckout rejects invalid quantity and under-minimum supplier", func(t *testing.T) {
		t.Parallel()
		invalidQuantity := &model.Cart{Suppliers: []model.Supplier{{OrderAmount: 1000, TotalAmount: 1000, ProductList: []model.Product{{Price: 1000, Quantity: 0}}}}}
		assert.False(t, checkCartForCheckout(invalidQuantity))

		belowMinimum := &model.Cart{Suppliers: []model.Supplier{{OrderAmount: 5000, TotalAmount: 3000, ProductList: []model.Product{{Price: 3000, Quantity: 1}}}}}
		assert.False(t, checkCartForCheckout(belowMinimum))
	})
}
