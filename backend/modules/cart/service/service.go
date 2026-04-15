package service

import (
	"context"
	"diploma/modules/cart/model"
	userModel "diploma/modules/user/model"
	"diploma/pkg/client/db"
	"diploma/pkg/service"
)

type cartServ struct {
	service.BaseService
	cartRepo        ICartRepository
	productService  IProductService
	supplierService ISupplierClient
	orderService    IOrderClient
	addressBook     IAddressBook
	PaymentClient   IPaymentClient
	redis           IRedis
	txManager       db.TxManager
}

func NewService(
	cartRepository ICartRepository,
	productService IProductService,
	supplierService ISupplierClient,
	OrderClient IOrderClient,
	addressBook IAddressBook,
	PaymentClient IPaymentClient,
	redis IRedis,
	txManager db.TxManager,
) *cartServ {
	return &cartServ{
		BaseService:     service.NewBaseService("cart"),
		cartRepo:        cartRepository,
		supplierService: supplierService,
		orderService:    OrderClient,
		addressBook:     addressBook,
		PaymentClient:   PaymentClient,
		redis:           redis,
		txManager:       txManager,
		productService:  productService,
	}
}

type IPaymentClient interface {
	PaymentRequest(orderID, amount, currency, description string) (model.CheckoutResponse, error)
}

type IProductService interface {
	ProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error)
}

type ISupplierClient interface {
	SupplierListByIDList(ctx context.Context, IDList []int64) ([]model.Supplier, error)
}

type IOrderClient interface {
	CreateOrder(ctx context.Context, cart *model.Cart, address model.DeliveryAddress, paymentID string) error
}

type IAddressBook interface {
	GetById(ctx context.Context, id int64) (*userModel.Address, error)
}

type IRedis interface {
	SavePaymentOrder(ctx context.Context, paymentOrder model.PaymentOrder) error
	PaymentOrder(ctx context.Context, orderID string) (model.PaymentOrder, error)
}

type ICartRepository interface {
	UpdateItemQuantity(ctx context.Context, cartId, productId, supplierId int64, quantity int) error
	ItemQuantity(ctx context.Context, cartId, productId, supplierId int64) (int, error)
	Cart(ctx context.Context, userID int64) (*model.Cart, error)
	CreateCart(ctx context.Context, userID int64) (int64, error)
	AddItem(ctx context.Context, input *model.PutCartQuery) error
	UpdateCartTotal(ctx context.Context, cartID int64, total int) error
	DeleteCart(ctx context.Context, cartID int64) error
	GetCartItems(ctx context.Context, cartID int64) ([]model.Supplier, error)
	DeleteCartItems(ctx context.Context, cartID int64) error
	DeleteItem(ctx context.Context, cartID, productId, supplierId int64) error
}
