package usecase

import (
	"context"

	"github.com/The-Ogulgozel/Banking-system/internal/models"
)

type AccountsUsecaseInterface interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id int64) (*models.Account, error)
	GetAll(ctx context.Context, params *models.ListAccounts) ([]*models.Account, int64, error)
	Delete(ctx context.Context, id int64) error
}
