package model

import "errors"

var (
	ErrUnauthorized = errors.New("api: unauthorized")
)

type SetAddressInput struct {
	Address Address `json:"address"`
}

type GetAddressResponse struct {
	AddressList []Address `json:"address_list"`
}

type Address struct {
	ID          int64  `json:"id"`
	Street      string `json:"street"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	Err string `json:"error"`
}
