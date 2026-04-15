package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	fondyCheckoutURLEnvName      = "FONDY_CHEKCOUT_URL"
	fondyMerchantIDEnvName       = "FONDY_MERCHANT_ID"
	fondyMerchantPasswordEnvName = "FONDY_MERCHANT_PASSWORD"
	fondyCallbackURLEnvName      = "FONDY_CALLBACK_URL"
)

type PaymentConfig interface {
	CheckoutURL() string
	MerchantID() string
	MerchantPassword() string
	CallbackURL() string
}

type paymentConfig struct {
	checkoutURL      string
	merchantID       string
	merchantPassword string
	callbackURL      string
}

func NewPaymentConfig() (PaymentConfig, error) {
	checkoutURL := os.Getenv(fondyCheckoutURLEnvName)
	if len(checkoutURL) == 0 {
		return nil, errors.New("missing FONDY_CHEKCOUT_URL")
	}
	merchantID := os.Getenv(fondyMerchantIDEnvName)
	if len(merchantID) == 0 {
		return nil, errors.New("missing FONDY_MERCHANT_ID")
	}
	merchantPassword := os.Getenv(fondyMerchantPasswordEnvName)
	if len(merchantPassword) == 0 {
		return nil, errors.New("missing FONDY_MERCHANT_PASSWORD")
	}

	callbackURL := os.Getenv(fondyCallbackURLEnvName)
	if len(callbackURL) == 0 {
		return nil, errors.New("missing FONDY_CALLBACK_URL")
	}

	fmt.Println(checkoutURL)
	return &paymentConfig{
		checkoutURL:      checkoutURL,
		merchantID:       merchantID,
		merchantPassword: merchantPassword,
		callbackURL:      callbackURL,
	}, nil
}

func (cfg *paymentConfig) CheckoutURL() string {
	return cfg.checkoutURL
}

func (cfg *paymentConfig) MerchantID() string {
	return cfg.merchantID
}

func (cfg *paymentConfig) MerchantPassword() string {
	return cfg.merchantPassword
}

func (cfg *paymentConfig) CallbackURL() string {
	return cfg.callbackURL
}
