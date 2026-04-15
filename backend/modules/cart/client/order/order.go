package order

import (
	"context"
	"diploma/modules/cart/model"
	orderModel "diploma/modules/order/model"
	order "diploma/modules/order/service"
	"fmt"
)

type OrderClient struct {
	orderService *order.OrderService
}

func NewClient(orderService *order.OrderService) *OrderClient {
	return &OrderClient{orderService: orderService}
}

func (a *OrderClient) CreateOrder(ctx context.Context, cart *model.Cart, address model.DeliveryAddress, paymentID string) error {
	orders := cartToOrder(cart, address, paymentID)
	fmt.Println(orders)
	if err := a.orderService.CreateOrder(ctx, orders); err != nil {
		return err
	}

	return nil
}

func cartToOrder(cart *model.Cart, address model.DeliveryAddress, paymentID string) []*orderModel.Order {
	res := make([]*orderModel.Order, 0, len(cart.Suppliers))
	for _, supplier := range cart.Suppliers {

		order := orderModel.Order{
			CustomerID:      cart.CustomerID,
			SupplierID:      supplier.ID,
			DeliveryAddress: address.Street,
			DeliveryComment: address.Description,
			PaymentID:       paymentID,
			PaymentStatus:   "paid",
			ProductList:     cartProductListToOrderProductList(supplier.ProductList),
		}
		res = append(res, &order)
	}
	return res
}

func cartProductListToOrderProductList(cartProducts []model.Product) []*orderModel.OrderProduct {
	orderProducts := make([]*orderModel.OrderProduct, 0, len(cartProducts))
	for _, cp := range cartProducts {
		orderProducts = append(orderProducts, &orderModel.OrderProduct{
			ProductID: cp.ID,
			Quantity:  cp.Quantity,
			Price:     cp.Price,
		})
	}
	return orderProducts
}
