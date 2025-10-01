package wallet

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Wallet struct {
	//Id       int       `json:"id"`
	ValletId uuid.UUID `json:"valletId"`
	Balance  float64   `json:"balance"`
}

type WalletTransactions struct {
	Id            int       `json:"id"`
	ValletId      uuid.UUID `json:"valletId" binding:"required"`
	OperationType string    `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64   `json:"amount" binding:"required"`
}
