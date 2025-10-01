package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/KatenkaKet/wallet"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type WalletPsql struct {
	db *sqlx.DB
}

func NewWalletPsql(db *sqlx.DB) *WalletPsql {
	return &WalletPsql{db: db}
}

func (w *WalletPsql) GetBalance(ctx context.Context, uid uuid.UUID) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var balance float64

	query := fmt.Sprintf("SELECT balance FROM %s WHERE ValletId=$1", walletTable)
	err := w.db.QueryRowContext(ctx, query, uid).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance for wallet %s: %w", uid.UUID.String(), err)
	}
	return balance, nil
}

func (w *WalletPsql) UpdateBalance(ctx context.Context, uid uuid.UUID, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	_, err = tx.ExecContext(ctx, fmt.Sprintf(
		`SELECT 1 FROM %s WHERE valletid = $1 FOR UPDATE`, walletTable), uid)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to lock wallet %s: %w", uid.UUID.String(), err)
	}

	_, err = tx.ExecContext(ctx, fmt.Sprintf(
		`UPDATE %s SET balance = balance + $1 WHERE valletid = $2`, walletTable), amount, uid)
	if err != nil {
		tx.Rollback()

		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23514" { // check_violation
			return fmt.Errorf("insufficient funds for wallet %s", uid.UUID.String())
		}

		return fmt.Errorf("failed to update balance for wallet %s: %w", uid.UUID.String(), err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx for wallet %s: %w", uid.UUID.String(), err)
	}

	return nil
}

func (w *WalletPsql) CreateTransaction(ctx context.Context, WT wallet.WalletTransactions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := fmt.Sprintf(`INSERT INTO %s (valletId, operation_type, amount) VALUES ($1, $2, $3)`, walletTRXTable)
	_, err := w.db.ExecContext(ctx, query, WT.ValletId, WT.OperationType, WT.Amount)
	if err != nil {
		return fmt.Errorf("failed to insert transaction for wallet %s: %w", WT.ValletId.UUID.String(), err)
	}
	return nil
}
