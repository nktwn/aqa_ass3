package model

type CheckoutRequest struct {
	AddressID int64 `json:"address_id" binding:"required"`
}

type CheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
	OrderID     string `json:"order_id"`
	PaymentID   string `json:"payment_id"`
}
