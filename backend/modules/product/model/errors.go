package model

import "errors"

var (
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateNumber = errors.New("models: duplicate email")

	ErrNoRows = errors.New("models: no rows")
)
