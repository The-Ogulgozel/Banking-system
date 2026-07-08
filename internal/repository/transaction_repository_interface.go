package repository

import (
	"context"

	"github.com/The-Ogulgozel/Banking-system/internal/models"
)

type TransactionRepoInterface interface {
	Deposit(ctx context.Context, id int64, amount int64) (*models.Transaction, error)
	Withdraw(ctx context.Context, id int64, amount int64) (*models.Transaction, error)
	Transfer(ctx context.Context, fromAccountID, toAccountID, amount int64) (*models.Transaction, error)
	GetByAccountId(ctx context.Context, param *models.ListTransactions) ([]*models.Transaction, int64, error)
}
