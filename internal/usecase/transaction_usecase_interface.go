package usecase

import (
	"context"

	"github.com/The-Ogulgozel/Banking-system/internal/models"
)

type TransactionUseCaseInterface interface {
	Deposit(ctx context.Context, id int64, amount int64) (*models.Transaction, error)
	Withdraw(ctx context.Context, id int64, amount int64) (*models.Transaction, error)
	Transfer(ctx context.Context, fromId, toId, amount int64) (*models.Transaction, error)
	ListByAccountID(ctx context.Context, param *models.ListTransactions) ([]*models.Transaction, int64, error)
}
