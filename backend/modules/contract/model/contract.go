package model

import (
	"database/sql"
	"time"
)

const (
	StatusCreated = iota
	StatusSignedBySupplier
	StatusSignedByCustomer
	StatusCompleted
)

type Contract struct {
	ID          int64
	OrderID     int64
	SupplierID  int64
	CustomerID  int64
	Content     string
	SupplierSig sql.NullString
	CustomerSig sql.NullString
	Status      int
	CreatedAt   time.Time
	SignedAt    sql.NullTime
}

type SignatureRequest struct {
	ContractID int64  `json:"contract_id" binding:"required"`
	Signature  string `json:"signature" binding:"required"`
}
