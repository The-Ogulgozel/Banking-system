package models

import "time"

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID          int64
	AccountID   int64
	ToAccountID *int64
	Type        TransactionType
	Amount      int64
	CreatedAt   time.Time
}

type ListTransactions struct {
	AccountID int64
	Type      *TransactionType
	Page      *int64
	Limit     *int64
	StartDate *time.Time
	EndDate   *time.Time
}
