package repository

import "github.com/jmoiron/sqlx"

type Wallet interface {
}

type Repository struct {
	Wallet
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
