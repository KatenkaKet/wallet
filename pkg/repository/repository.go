package repository

import (
	"context"

	"github.com/KatenkaKet/wallet"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jmoiron/sqlx"
)

type Wallet interface {
	GetBalance(ctx context.Context, uuid uuid.UUID) (float64, error)
	UpdateBalance(ctx context.Context, uuid uuid.UUID, amount float64) error
	CreateTransaction(ctx context.Context, WT wallet.WalletTransactions) error
}

type Repository struct {
	Wallet
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Wallet: NewWalletPsql(db),
	}
}
