package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/KatenkaKet/wallet"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func uuidFromString(s string) uuid.UUID {
	var id uuid.UUID
	if err := id.DecodeText(nil, []byte(s)); err != nil {
		log.Fatal(err)
	}
	return id
}

func TestWalletPsql_GetBalance(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewWalletPsql(db)
	uid := uuidFromString("11111111-1111-1111-1111-111111111111")

	testTable := []struct {
		name            string
		mockSetup       func()
		expectedBalance float64
		expectError     bool
	}{
		{
			name: "success",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(100.5)
				mock.ExpectQuery(fmt.Sprintf(`SELECT balance FROM %s WHERE ValletId=\$1`, walletTable)).
					WithArgs(uid).
					WillReturnRows(rows)
			},
			expectedBalance: 100.5,
			expectError:     false,
		},
		{
			name: "wallet not found",
			mockSetup: func() {
				mock.ExpectQuery(fmt.Sprintf(`SELECT balance FROM %s WHERE ValletId=\$1`, walletTable)).
					WithArgs(uid).
					WillReturnError(sql.ErrNoRows)
			},
			expectedBalance: 0.0,
			expectError:     true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			balance, err := r.GetBalance(context.Background(), uid)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedBalance, balance)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestWalletPsql_CreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	w := NewWalletPsql(db)
	uid := uuidFromString("11111111-1111-1111-1111-111111111111")

	testTable := []struct {
		name      string
		inputWT   wallet.WalletTransactions
		mockSetup func()
		expectErr bool
	}{
		{
			name: "success",
			inputWT: wallet.WalletTransactions{
				ValletId:      uid,
				OperationType: "DEPOSIT",
				Amount:        100.5,
			},
			mockSetup: func() {
				mock.ExpectExec(
					fmt.Sprintf(`INSERT INTO %s \(valletId, operation_type, amount\) VALUES \(\$1, \$2, \$3\)`, walletTRXTable)).
					WithArgs(uid, "DEPOSIT", 100.5).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectErr: false,
		},
		{
			name: "insert error",
			inputWT: wallet.WalletTransactions{
				ValletId:      uid,
				OperationType: "WITHDRAW",
				Amount:        50,
			},
			mockSetup: func() {
				mock.ExpectExec(
					fmt.Sprintf(`INSERT INTO %s \(valletId, operation_type, amount\) VALUES \(\$1, \$2, \$3\)`, walletTRXTable)).
					WithArgs(uid, "WITHDRAW", 50).
					WillReturnError(errors.New("insert failed"))
			},
			expectErr: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			err := w.CreateTransaction(context.Background(), test.inputWT)

			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestWalletPsql_UpdateBalance(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatalf("unexpected error opening stub db: %s", err)
	}
	defer db.Close()

	w := NewWalletPsql(db)
	uid := uuidFromString("11111111-1111-1111-1111-111111111111")

	testTable := []struct {
		name      string
		amount    float64
		mockSetup func()
		expectErr bool
	}{
		{
			name:   "success",
			amount: 100.0,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(fmt.Sprintf(`SELECT 1 FROM %s WHERE valletid = \$1 FOR UPDATE`, walletTable)).
					WithArgs(uid.UUID.String()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(fmt.Sprintf(`UPDATE %s SET balance = balance \+ \$1 WHERE valletid = \$2`, walletTable)).
					WithArgs(100.0, uid.UUID.String()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name:   "lock error",
			amount: 50.0,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(fmt.Sprintf(`SELECT 1 FROM %s WHERE valletid = \$1 FOR UPDATE`, walletTable)).
					WithArgs(uid.UUID.String()).
					WillReturnError(errors.New("lock failed"))
				mock.ExpectRollback()
			},
			expectErr: true,
		},
		{
			name:   "insufficient funds",
			amount: -200.0,
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec(fmt.Sprintf(`SELECT 1 FROM %s WHERE valletid = \$1 FOR UPDATE`, walletTable)).
					WithArgs(uid.UUID.String()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// эмулируем ошибку postgres check constraint violation (23514)
				mock.ExpectExec(fmt.Sprintf(`UPDATE %s SET balance = balance \+ \$1 WHERE valletid = \$2`, walletTable)).
					WithArgs(-200.0, uid.UUID.String()).
					WillReturnError(&pq.Error{Code: "23514"})

				mock.ExpectRollback()
			},
			expectErr: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockSetup()

			err := w.UpdateBalance(context.Background(), uid, test.amount)

			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// проверяем, что моковые ожидания выполнены
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
