package converter

import (
	"diploma/modules/cart/model"
	"fmt"
)

func ToServiceCheckoutFromClient(checkout map[string]string) (model.CheckoutResponse, error) {

	checkoutURL, ok := checkout["checkout_url"]
	if !ok {
		return model.CheckoutResponse{}, fmt.Errorf("checkout_url not found in response")
	}
	paymentID, ok := checkout["payment_id"]
	if !ok {
		return model.CheckoutResponse{}, fmt.Errorf("payment_id not found in response")
	}
	responseStatus, ok := checkout["response_status"]
	if !ok {
		return model.CheckoutResponse{}, fmt.Errorf("response_status not found in response")
	}

	return model.CheckoutResponse{
		CheckoutURL:    checkoutURL,
		PaymentID:      paymentID,
		ResponseStatus: responseStatus,
	}, nil
}
