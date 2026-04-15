package repository

import (
	"diploma/pkg/client/db"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}
