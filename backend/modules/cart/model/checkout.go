package model

const (
	PaymentStatusApproved = "approved"
	PaymentStatusFailed  = "failed"
)

type CheckoutResponse struct {
	CheckoutURL    string
	PaymentID      string
	OrderID        string
	ResponseStatus string
}

type CheckoutInput struct {
	AddressID int64
}

type DeliveryAddress struct {
	ID          int64
	Street      string
	Description string
}

type PaymentOrder struct {
	ID              string
	Cart            Cart
	DeliveryAddress DeliveryAddress
	CheckoutResponse
}

type CommitCheckout struct {
	OrderID       string
	PaymentStatus string
}
