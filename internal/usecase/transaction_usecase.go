package usecase

import (
	"context"

	appError "github.com/The-Ogulgozel/Banking-system/internal/errors"
	"github.com/The-Ogulgozel/Banking-system/internal/models"
	"github.com/The-Ogulgozel/Banking-system/internal/repository"
)

type TransactionUsecase struct {
	transactionRepo repository.TransactionRepoInterface
	accountRepo     repository.AccountRepositoryInterface
}

func NewTransactionUsecase(transactionRepo repository.TransactionRepoInterface, accountRepo repository.AccountRepositoryInterface) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (u *TransactionUsecase) Deposit(ctx context.Context, id int64, amount int64) (*models.Transaction, error) {
	account, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if account.IsLocked {
		return nil, appError.ErrAccountLocked
	}

	if account.Balance+amount > 1_000_000_000_000 {
		return nil, appError.ErrBalanceLimitExceeded
	}

	t, err := u.transactionRepo.Deposit(ctx, id, amount)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *TransactionUsecase) Withdraw(ctx context.Context, id int64, amount int64) (*models.Transaction, error) {
	account, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if account.IsLocked {
		return nil, appError.ErrAccountLocked
	}

	if account.Balance < amount {
		return nil, appError.ErrNotEnoughBalance
	}

	t, err := u.transactionRepo.Withdraw(ctx, id, amount)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *TransactionUsecase) Transfer(ctx context.Context, fromId, toId, amount int64) (*models.Transaction, error) {
	fromAccount, err := u.accountRepo.GetByID(ctx, fromId)
	if err != nil {
		return nil, err
	}

	if fromAccount.IsLocked {
		return nil, appError.ErrAccountLocked
	}

	if fromAccount.Balance < amount {
		return nil, appError.ErrNotEnoughBalance
	}

	toAccount, err := u.accountRepo.GetByID(ctx, toId)
	if err != nil {
		return nil, err
	}

	if toAccount.IsLocked {
		return nil, appError.ErrAccountLocked
	}

	if toAccount.Balance+amount > 1_000_000_000_000 {
		return nil, appError.ErrBalanceLimitExceeded
	}

	if fromAccount.Currency != toAccount.Currency {
		return nil, appError.ErrCurrencyMismatch
	}

	t, err := u.transactionRepo.Transfer(ctx, fromId, toId, amount)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *TransactionUsecase) ListByAccountID(ctx context.Context, params *models.ListTransactions) ([]*models.Transaction, int64, error) {
	if params.Page == nil {
		var i int64 = 1
		params.Page = &i
	}
	if params.Limit == nil {
		var i int64 = 10
		params.Limit = &i
	}

	return u.transactionRepo.GetByAccountId(ctx, params)
}
