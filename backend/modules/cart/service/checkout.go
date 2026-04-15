package service

import (
	"context"
	"diploma/modules/cart/model"
	"fmt"
	"github.com/google/uuid"
	"strconv"
)

func (s *cartServ) Checkout(ctx context.Context, userID int64, input model.CheckoutInput) (model.CheckoutResponse, error) {
	var checkout model.CheckoutResponse
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		address, errTx := s.addressBook.GetById(ctx, input.AddressID)
		if errTx != nil {
			return errTx
		}
		if address.UserID != userID {
			return fmt.Errorf("address does not belong to current user")
		}

		// get cart
		cart, errTx := s.Cart(ctx, userID)
		if errTx != nil {
			return errTx
		}
		cart.CustomerID = userID
		if !checkCartForCheckout(cart) {
			return model.ErrInvalidCart
		}

		// process payment
		amount := getCartAmount(cart)

		paymentOrderID := generateOrderID(cart.CustomerID, cart.ID)
		checkout, errTx = s.PaymentClient.PaymentRequest(paymentOrderID, strconv.Itoa(amount)+"00", "USD", "Payment for order")
		if errTx != nil {
			return errTx
		}
		checkout.OrderID = paymentOrderID
		paymentOrder := model.PaymentOrder{
			ID: paymentOrderID,
			CheckoutResponse: checkout,
			Cart: *cart,
			DeliveryAddress: model.DeliveryAddress{
				ID:          address.ID,
				Street:      address.Street,
				Description: address.Description,
			},
		}

		// save payment order to redis
		return s.redis.SavePaymentOrder(ctx, paymentOrder)

	})

	if err != nil {
		return model.CheckoutResponse{}, err
	}
	return checkout, nil
}

func (s *cartServ) CommitCheckout(ctx context.Context, commitCheckout model.CommitCheckout) error {
	fmt.Println("CommitCheckout", commitCheckout)
	switch commitCheckout.PaymentStatus {
	case model.PaymentStatusApproved:
		paymentOrder, err := s.redis.PaymentOrder(ctx, commitCheckout.OrderID)
		if err != nil {
			return err
		}
		err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			errTx := s.orderService.CreateOrder(ctx, &paymentOrder.Cart, paymentOrder.DeliveryAddress, paymentOrder.PaymentID)
			if errTx != nil {
				return errTx
			}

			errTx = s.cartRepo.DeleteCartItems(ctx, paymentOrder.Cart.ID)
			if errTx != nil {
				return errTx
			}
			errTx = s.cartRepo.DeleteCart(ctx, paymentOrder.Cart.ID)
			if errTx != nil {
				return errTx
			}
			return nil
		})
		if err != nil {
			return err
		}
	case model.PaymentStatusFailed:

	}
	return nil
}

func generateOrderID(userID, cartID int64) string {
	id := uuid.New().String()
	return fmt.Sprintf("%d-%d-%s", userID, cartID, id)
}

// func (s *cartServ) CommitCheckout(ctx context.Context, userID int64) (*model.Cart, error) {

// }

func checkCartForCheckout(cart *model.Cart) bool {
	for _, supplier := range cart.Suppliers {
		sum := 0
		for _, product := range supplier.ProductList {
			if product.Quantity <= 0 || product.Price <= 0 {
				return false
			}
			sum += product.Price * product.Quantity
		}

		if supplier.TotalAmount < supplier.OrderAmount {
			return false
		}
	}
	return true
}

func getCartAmount(cart *model.Cart) int {
	totalSum := 0
	for _, supplier := range cart.Suppliers {
		totalSum += supplier.TotalAmount
	}
	return totalSum
}
