package wallet

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Wallet struct {
	Id       int       `json : "id"`
	Valletid uuid.UUID `json :"valletid"`
	Ballance float64   `json :"ballance"`
}

type WalletTransactions struct {
	Id            int       `json : "id"`
	ValletId      uuid.UUID `json:"valletId"`
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"`
}
