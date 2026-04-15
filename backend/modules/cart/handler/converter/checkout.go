package converter

import (
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
	"errors"
)

func ToAPICheckoutFromService(checkoutResponse model.CheckoutResponse) modelApi.CheckoutResponse {
	return modelApi.CheckoutResponse{
		CheckoutURL: checkoutResponse.CheckoutURL,
		OrderID:     checkoutResponse.OrderID,
		PaymentID:   checkoutResponse.PaymentID,
	}
}

func ToServiceCheckoutFromApi(data map[string]interface{}) (model.CommitCheckout, error) {
	orderID, ok := data["order_id"].(string)
	if !ok {
		return model.CommitCheckout{}, errors.New("order_id not found in response")
	}
	paymentStatus, ok := data["order_status"].(string)
	if !ok {
		return model.CommitCheckout{}, errors.New("order_status not found in response")
	}

	return model.CommitCheckout{
		OrderID:       orderID,
		PaymentStatus: paymentStatus,
	}, nil
}
