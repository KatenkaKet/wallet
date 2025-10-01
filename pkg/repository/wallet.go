package repository

import (
	"context"
	"fmt"

	"github.com/KatenkaKet/wallet"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jmoiron/sqlx"
)

type WalletPsql struct {
	db *sqlx.DB
}

func NewWalletPsql(db *sqlx.DB) *WalletPsql {
	return &WalletPsql{db: db}
}

func (w *WalletPsql) GetBalance(ctx context.Context, uid uuid.UUID) (float64, error) {
	var balance float64

	query := fmt.Sprintf("SELECT balance FROM %s WHERE ValletId=$1", walletTable)
	err := w.db.QueryRowContext(ctx, query, uid).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance for wallet %s: %w", uid.UUID.String(), err)
	}
	return balance, nil
}

func (w *WalletPsql) UpdateBalance(ctx context.Context, uid uuid.UUID, amount float64) error {
	query := fmt.Sprintf(`UPDATE %s SET balance = balance + $1 WHERE valletid = $2`, walletTable)
	_, err := w.db.ExecContext(ctx, query, amount, uid)
	if err != nil {
		return fmt.Errorf("failed to update balance for wallet %s: %w", uid.UUID.String(), err)
	}
	return nil
}

func (w *WalletPsql) CreateTransaction(ctx context.Context, WT wallet.WalletTransactions) error {
	query := fmt.Sprintf(`INSERT INTO %s (valletId, operation_type, amount) VALUES ($1, $2, $3)`, walletTRXTable)
	_, err := w.db.ExecContext(ctx, query, WT.ValletId, WT.OperationType, WT.Amount)
	if err != nil {
		return fmt.Errorf("failed to insert transaction for wallet %s: %w", WT.ValletId.UUID.String(), err)
	}
	return nil
}

//func (w *WalletPsql) CreateTransaction(walletTRX wallet.WalletTransactions) error {
//	//query :=
//}
