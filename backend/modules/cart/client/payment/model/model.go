package model

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/structs"
)

type APIResponse struct {
	Response map[string]string `json:"response"`
}

type APIRequest struct {
	Request interface{} `json:"request"`
}

type CheckoutRequest struct {
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	OrderDesc         string `json:"order_desc"`
	Signature         string `json:"signature"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	ResponseURL       string `json:"response_url"`
	ServerCallbackURL string `json:"server_callback_url"`
	SenderEmail       string `json:"sender_email"`
	Language          string `json:"language"`
	ProductID         string `json:"product_id"`
}

func (r *CheckoutRequest) SetSignature(password string) {

	params := structs.Map(r)

	var key []string
	for k := range params {
		key = append(key, k)
	}
	sort.Strings(key)
	values := []string{}

	for _, v := range key {
		value := params[v].(string)
		if value != "" {
			values = append(values, value)
		}
	}

	r.Signature = generateSignature(password, values)
}

func generateSignature(password string, values []string) string {
	newValues := []string{password}
	newValues = append(newValues, values...)
	signatureString := strings.Join(newValues, "|")
	fmt.Println("Signature String:", signatureString)
	hash := sha1.New()
	hash.Write([]byte(signatureString))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
