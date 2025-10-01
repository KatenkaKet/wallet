package service

import (
	"context"

	"github.com/KatenkaKet/wallet"
	"github.com/KatenkaKet/wallet/pkg/repository"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Wallet interface {
	GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error)
	UpdateBalance(ctx context.Context, WT wallet.WalletTransactions) error
}

type Service struct {
	Wallet
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Wallet: NewWalletService(repo.Wallet),
	}
}
