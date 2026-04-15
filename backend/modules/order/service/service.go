package service

import (
	"context"
	"diploma/modules/order/model"
	"diploma/pkg/client/db"
	"diploma/pkg/service"
)

type OrderService struct {
	service.BaseService
	orderRepo      IOrderRepository
	supplierClient ISupplierClient
	productClient  IProductClient
	contractClient IContractService
	notificationService INotificationService
	txManager      db.TxManager
}

func NewService(
	repo IOrderRepository,
	supplierClient ISupplierClient,
	productClient IProductClient,
	contractClient IContractService,
	notificationService INotificationService,
	tx db.TxManager,
) *OrderService {
	return &OrderService{
		BaseService:         service.NewBaseService("order"),
		orderRepo:           repo,
		supplierClient:      supplierClient,
		productClient:       productClient,
		contractClient:      contractClient,
		notificationService: notificationService,
		txManager:           tx,
	}
}

type IOrderRepository interface {
	ICreateOrderRepo
	IOrderRepo
	UpdateOrderStatus(ctx context.Context, orderID int64, newStatus int) error
	GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error)
}

type ISupplierClient interface {
	Supplier(ctx context.Context, id int64) (*model.Supplier, error)
}

type IProductClient interface {
	Product(ctx context.Context, id int64) (*model.Product, error)
}

type IContractService interface {
	CreateContract(ctx context.Context, orderID, supplierID, customerID int64, content string) (int64, error)
}

type INotificationService interface {
	Notify(ctx context.Context, userID int64, title, message, kind string, orderID *int64) error
}
