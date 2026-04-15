package payment

import (
	"diploma/modules/cart/client/payment/converter"
	clientModel "diploma/modules/cart/client/payment/model"
	"diploma/modules/cart/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type PaymentClient struct {
	checkoutURL      string
	merchantID       string
	merchantPassword string
	callbackURL      string
}

func NewPaymentClient(checkoutURL, merchantID, merchantPassword, callbackURL string) *PaymentClient {
	return &PaymentClient{
		checkoutURL:      checkoutURL,
		merchantID:       merchantID,
		merchantPassword: merchantPassword,
		callbackURL:      callbackURL,
	}
}

func (p *PaymentClient) PaymentRequest(orderID, amount, currency, orderDesc string) (model.CheckoutResponse, error) {
	checkoutRequest := clientModel.CheckoutRequest{
		OrderID:           orderID,
		MerchantID:        p.merchantID,
		OrderDesc:         orderDesc,
		Amount:            amount,
		Currency:          currency,
		ServerCallbackURL: p.callbackURL,
	}

	checkoutRequest.SetSignature(p.merchantPassword)

	apiRequest := clientModel.APIRequest{
		Request: checkoutRequest,
	}

	requestBody, err := json.Marshal(apiRequest)
	if err != nil {
		return model.CheckoutResponse{}, fmt.Errorf("error encoding request: %w", err)
	}

	resp, err := http.Post(p.checkoutURL, "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		return p.mockCheckoutResponse(orderID), nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return p.mockCheckoutResponse(orderID), nil
	}
	if resp.StatusCode != http.StatusOK {
		return p.mockCheckoutResponse(orderID), nil
	}
	var apiResponse clientModel.APIResponse
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return p.mockCheckoutResponse(orderID), nil
	}

	return converter.ToServiceCheckoutFromClient(apiResponse.Response)
}

func (p *PaymentClient) mockCheckoutResponse(orderID string) model.CheckoutResponse {
	mockURL := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   "/api/cart/mock-success",
	}

	query := mockURL.Query()
	query.Set("order_id", orderID)
	mockURL.RawQuery = query.Encode()

	return model.CheckoutResponse{
		CheckoutURL:    mockURL.String(),
		PaymentID:      orderID,
		ResponseStatus: "mock",
	}
}
