package handler

import (
	"context"
	"diploma/modules/cart/model"
)

type CartHandler struct {
	service ICartService
}

func NewHandler(service ICartService) *CartHandler {
	return &CartHandler{service: service}
}

type ICartService interface {
	ClearCart(ctx context.Context, userID int64) error
	Checkout(ctx context.Context, userID int64, input model.CheckoutInput) (model.CheckoutResponse, error)
	Cart(ctx context.Context, userID int64) (*model.Cart, error)
	AddProductToCard(ctx context.Context, input *model.PutCartQuery) error
	DeleteProductFromCart(ctx context.Context, input *model.PutCartQuery) error
	CommitCheckout(ctx context.Context, commitCheckout model.CommitCheckout) error
	// DeleteProductFromCart(ctx context.Context, input *model.DeleteProductQuery) error
}
