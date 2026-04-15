package repository

import "diploma/pkg/client/db"

type OrderRepo struct {
	db db.Client
}

const (
	// ======== supplier table ========

)

func NewRepository(db db.Client) *OrderRepo {
	return &OrderRepo{db: db}
}
