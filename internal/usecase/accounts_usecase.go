package usecase

import (
	"context"

	appError "github.com/The-Ogulgozel/Banking-system/internal/errors"
	"github.com/The-Ogulgozel/Banking-system/internal/models"
	"github.com/The-Ogulgozel/Banking-system/internal/repository"
)

type AccountUsecase struct {
	AccountRepo repository.AccountRepositoryInterface
}

func NewAccountUsecase(accountRepo repository.AccountRepositoryInterface) *AccountUsecase {
	return &AccountUsecase{
		AccountRepo: accountRepo,
	}
}

func (u *AccountUsecase) Create(ctx context.Context, accounts *models.Account) error {
	return u.AccountRepo.Create(ctx, accounts)
}

func (u *AccountUsecase) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	return u.AccountRepo.GetByID(ctx, id)
}

func (u *AccountUsecase) GetAll(ctx context.Context, params *models.ListAccounts) ([]*models.Account, int64, error) {
	if params.Page == nil {
		var i int64 = 1
		params.Page = &i
	}
	if params.Limit == nil {
		var i int64 = 10
		params.Limit = &i
	}
	return u.AccountRepo.GetAll(ctx, params)
}

func (u *AccountUsecase) Delete(ctx context.Context, id int64) error {
	account, err := u.AccountRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if account.IsLocked {
		return appError.ErrAccountLocked
	}
	if account.Balance != 0 {
		return appError.ErrAccountBalanceNotEmpty
	}

	return u.AccountRepo.Delete(ctx, id)
}
