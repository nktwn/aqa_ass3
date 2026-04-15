package service

import (
	"context"
	"diploma/modules/order/model"
	"diploma/pkg/client/db"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockOrderRepository struct{ mock.Mock }

func (m *mockOrderRepository) CreateOrder(ctx context.Context, order *model.Order) (int64, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockOrderRepository) CreateOrderProduct(ctx context.Context, orderProduct *model.OrderProduct) error {
	return m.Called(ctx, orderProduct).Error(0)
}
func (m *mockOrderRepository) OrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*model.Order), args.Error(1)
}
func (m *mockOrderRepository) OrdersBySupplierID(ctx context.Context, supplierID int64) ([]*model.Order, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*model.Order), args.Error(1)
}
func (m *mockOrderRepository) OrderProducts(ctx context.Context, orderID int64) ([]*model.OrderProduct, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*model.OrderProduct), args.Error(1)
}
func (m *mockOrderRepository) UpdateOrderStatus(ctx context.Context, orderID int64, newStatus int) error {
	return m.Called(ctx, orderID, newStatus).Error(0)
}
func (m *mockOrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Order), args.Error(1)
}

type mockSupplierClient struct{ mock.Mock }
func (m *mockSupplierClient) Supplier(ctx context.Context, id int64) (*model.Supplier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Supplier), args.Error(1)
}

type mockProductClient struct{ mock.Mock }
func (m *mockProductClient) Product(ctx context.Context, id int64) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*model.Product), args.Error(1)
}

type mockContractService struct{ mock.Mock }
func (m *mockContractService) CreateContract(ctx context.Context, orderID, supplierID, customerID int64, content string) (int64, error) {
	args := m.Called(ctx, orderID, supplierID, customerID, content)
	return args.Get(0).(int64), args.Error(1)
}

type mockNotificationService struct{ mock.Mock }
func (m *mockNotificationService) Notify(ctx context.Context, userID int64, title, message, kind string, orderID *int64) error {
	return m.Called(ctx, userID, title, message, kind, orderID).Error(0)
}

type mockTxManager struct{ mock.Mock }
func (m *mockTxManager) ReadCommitted(ctx context.Context, h db.Handler) error {
	args := m.Called(ctx, h)
	if h != nil {
		if err := h(ctx); err != nil { return err }
	}
	return args.Error(0)
}

func TestUpdateOrderStatusBySupplier_CreatesContractForPendingOrder(t *testing.T) {
	t.Parallel()

	repo := new(mockOrderRepository)
	suppliers := new(mockSupplierClient)
	products := new(mockProductClient)
	contracts := new(mockContractService)
	notifications := new(mockNotificationService)
	tx := new(mockTxManager)
	service := NewService(repo, suppliers, products, contracts, notifications, tx)

	order := &model.Order{ID: 50, CustomerID: 10, SupplierID: 22, StatusID: model.Pending}
	tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("GetOrderByID", mock.Anything, int64(50)).Return(order, nil).Once()
	repo.On("UpdateOrderStatus", mock.Anything, int64(50), model.InProgress).Return(nil).Once()
	contracts.On("CreateContract", mock.Anything, int64(50), int64(22), int64(10), "Контракт для заказа #50").Return(int64(1001), nil).Once()
	notifications.On("Notify", mock.Anything, int64(10), "Статус заказа обновлён", mock.AnythingOfType("string"), "order_status", mock.Anything).Return(nil).Once()

	err := service.UpdateOrderStatusBySupplier(context.Background(), 22, 50, model.InProgress)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	contracts.AssertExpectations(t)
}

func TestUpdateOrderStatusBySupplier_RejectsUnauthorizedSupplier(t *testing.T) {
	t.Parallel()

	repo := new(mockOrderRepository)
	service := NewService(repo, new(mockSupplierClient), new(mockProductClient), new(mockContractService), new(mockNotificationService), new(mockTxManager))
	tx := service.txManager.(*mockTxManager)
	tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("GetOrderByID", mock.Anything, int64(50)).Return(&model.Order{ID: 50, SupplierID: 99, StatusID: model.Pending}, nil).Once()

	err := service.UpdateOrderStatusBySupplier(context.Background(), 22, 50, model.InProgress)

	assert.EqualError(t, err, "supplier 22 does not own order 50")
	repo.AssertNotCalled(t, "UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything)
}

func TestCancelOrderByCustomer_RejectsNonPendingOrder(t *testing.T) {
	t.Parallel()

	repo := new(mockOrderRepository)
	service := NewService(repo, new(mockSupplierClient), new(mockProductClient), new(mockContractService), new(mockNotificationService), new(mockTxManager))
	tx := service.txManager.(*mockTxManager)
	tx.On("ReadCommitted", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("GetOrderByID", mock.Anything, int64(60)).Return(&model.Order{ID: 60, CustomerID: 10, SupplierID: 22, StatusID: model.InProgress}, nil).Once()

	err := service.CancelOrderByCustomer(context.Background(), 10, 60)

	assert.EqualError(t, err, "only orders in Pending status can be cancelled")
	repo.AssertNotCalled(t, "UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything)
}
