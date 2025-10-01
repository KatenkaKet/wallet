package service

import (
	"context"

	"github.com/KatenkaKet/wallet"
	"github.com/KatenkaKet/wallet/pkg/repository"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type WalletService struct {
	repo repository.Wallet
}

func NewWalletService(repo repository.Wallet) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	return s.repo.GetBalance(ctx, walletID)
}

func (s *WalletService) UpdateBalance(ctx context.Context, WT wallet.WalletTransactions) error {
	if WT.OperationType == "WITHDRAW" {
		WT.Amount *= -1
	}

	err := s.repo.UpdateBalance(ctx, WT.ValletId, WT.Amount)
	if err != nil {
		return err
	}

	if WT.OperationType == "WITHDRAW" {
		WT.Amount *= -1
	}

	err = s.repo.CreateTransaction(ctx, WT)
	if err != nil {
		return err
	}

	return nil
}
