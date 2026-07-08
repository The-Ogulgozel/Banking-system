package repository

import (
	"context"
	"errors"
	"log"

	appError "github.com/The-Ogulgozel/Banking-system/internal/errors"
	"github.com/The-Ogulgozel/Banking-system/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepo(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Deposit(ctx context.Context, id int64, amount int64) (*models.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin tx: %v", err)
		return nil, appError.ErrInternalServer
	}
	defer tx.Rollback(ctx)

	if err := r.lockAccount(ctx, tx, id); err != nil {
		log.Printf("failed to lock account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}

	var balance int64
	if err := tx.QueryRow(ctx, `SELECT balance FROM accounts WHERE id=$1 FOR UPDATE`, id).Scan(&balance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound
		}
		log.Printf("failed to get balance for account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}

	newBalance := balance + amount
	if newBalance > 1_000_000_000_000 {
		return nil, appError.ErrBalanceLimitExceeded
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id=$2`,
		amount, id,
	)
	if err != nil {
		log.Printf("failed to update balance for account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}

	var t models.Transaction
	var typeStr string
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (account_id, amount, transaction_type)
		 VALUES ($1, $2, 'deposit')
		 RETURNING id, account_id, amount, transaction_type, created_at`,
		id, amount,
	).Scan(&t.ID, &t.AccountID, &t.Amount, &typeStr, &t.CreatedAt)
	if err != nil {
		log.Printf("failed to create transaction for account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}
	t.Type = models.TransactionType(typeStr)

	if err := r.unlockAccount(ctx, tx, id); err != nil {
		log.Printf("failed to unlock account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("failed to commit tx: %v", err)
		return nil, appError.ErrInternalServer
	}

	return &t, nil
}

func (r *TransactionRepository) Withdraw(ctx context.Context, id int64, amount int64) (*models.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin tx: %v", err)
		return nil, appError.ErrInternalServer
	}
	defer tx.Rollback(ctx)

	if err := r.lockAccount(ctx, tx, id); err != nil {
		log.Printf("failed to lock account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}
	var balance int64
	err = tx.QueryRow(ctx, `SELECT balance FROM accounts WHERE id=$1 FOR UPDATE`, id).Scan(&balance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound
		}
		log.Printf("failed to get balance for account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}
	if balance < amount {
		return nil, appError.ErrNotEnoughBalance
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance - $1 WHERE id=$2`,
		amount, id,
	)
	if err != nil {
		log.Printf("failed to update balance for account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}

	var t models.Transaction
	var typeStr string
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (account_id, amount, transaction_type)
		 VALUES ($1, $2, 'withdraw')
		 RETURNING id, account_id, amount, transaction_type, created_at`,
		id, amount,
	).Scan(&t.ID, &t.AccountID, &t.Amount, &typeStr, &t.CreatedAt)
	if err != nil {
		log.Printf("failed to create transaction for account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}
	t.Type = models.TransactionType(typeStr)

	if err := r.unlockAccount(ctx, tx, id); err != nil {
		log.Printf("failed to unlock account %d: %v", id, err)
		return nil, appError.ErrInternalServer
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("failed to commit tx: %v", err)
		return nil, appError.ErrInternalServer
	}

	return &t, nil
}

func (r *TransactionRepository) Transfer(ctx context.Context, fromId, toId, amount int64) (*models.Transaction, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin tx: %v", err)
		return nil, appError.ErrInternalServer
	}
	defer tx.Rollback(ctx)

	if err := r.lockAccount(ctx, tx, fromId); err != nil {
		log.Printf("failed to lock account %d: %v", fromId, err)
		return nil, appError.ErrInternalServer
	}
	if err := r.lockAccount(ctx, tx, toId); err != nil {
		log.Printf("failed to lock account %d: %v", toId, err)
		return nil, appError.ErrInternalServer
	}
	var fromBalance, toBalance int64
	var fromCurrency, toCurrency string
	if err := tx.QueryRow(ctx, `SELECT balance, currency FROM accounts WHERE id=$1 FOR UPDATE`, fromId).Scan(&fromBalance, &fromCurrency); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound
		}
		log.Printf("failed to fetch account %d: %v", fromId, err)
		return nil, appError.ErrInternalServer
	}

	if err := tx.QueryRow(ctx, `SELECT balance, currency FROM accounts WHERE id=$1 FOR UPDATE`, toId).Scan(&toBalance, &toCurrency); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound
		}
		log.Printf("failed to fetch account %d: %v", toId, err)
		return nil, appError.ErrInternalServer
	}

	if fromCurrency != toCurrency {
		return nil, appError.ErrCurrencyMismatch
	}
	if fromBalance < amount {
		return nil, appError.ErrNotEnoughBalance
	}
	newToBalance := toBalance + amount
	if newToBalance > 1_000_000_000_000 {
		return nil, appError.ErrBalanceLimitExceeded
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance - $1 WHERE id=$2`,
		amount, fromId,
	)
	if err != nil {
		log.Printf("failed to update balance for account %d: %v", fromId, err)
		return nil, appError.ErrInternalServer
	}

	_, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1 WHERE id=$2`,
		amount, toId,
	)
	if err != nil {
		log.Printf("failed to update balance for account %d: %v", toId, err)
		return nil, appError.ErrInternalServer
	}

	var t models.Transaction
	var typeStr string
	err = tx.QueryRow(ctx,
		`INSERT INTO transactions (account_id, to_account_id, amount, transaction_type)
		 VALUES ($1, $2, $3, 'transfer')
		 RETURNING id, account_id, to_account_id, amount, transaction_type, created_at`,
		fromId, toId, amount,
	).Scan(&t.ID, &t.AccountID, &t.ToAccountID, &t.Amount, &typeStr, &t.CreatedAt)
	if err != nil {
		log.Printf("failed to create transaction for account %d: %v", fromId, err)
		return nil, appError.ErrInternalServer
	}
	t.Type = models.TransactionType(typeStr)

	if err := r.unlockAccount(ctx, tx, fromId); err != nil {
		log.Printf("failed to unlock account %d: %v", fromId, err)
		return nil, appError.ErrInternalServer
	}
	if err := r.unlockAccount(ctx, tx, toId); err != nil {
		log.Printf("failed to unlock account %d: %v", toId, err)
		return nil, appError.ErrInternalServer
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("failed to commit tx: %v", err)
		return nil, appError.ErrInternalServer
	}

	return &t, nil
}

func (r *TransactionRepository) GetByAccountId(ctx context.Context, param *models.ListTransactions) ([]*models.Transaction, int64, error) {
	offset := (*param.Page - 1) * *param.Limit

	var typeFilter *string
	if param.Type != nil {
		t := string(*param.Type)
		typeFilter = &t
	}
	var total int64

	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM transactions
	WHERE (account_id=$1 OR to_account_id=$1)
	AND ($2::text IS NULL OR transaction_type::text = $2)
	AND ($3::date IS NULL OR created_at::date >= $3)
	AND ($4::date IS NULL OR created_at::date <= $4)`,
		param.AccountID, typeFilter, param.StartDate, param.EndDate,
	).Scan(&total)
	if err != nil {
		log.Printf("failed to get count of transactions for account %d: %v", param.AccountID, err)
		return nil, 0, appError.ErrInternalServer
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, account_id, to_account_id, amount, transaction_type, created_at 
	FROM transactions
	WHERE (account_id=$1 OR to_account_id=$1)
	AND ($2::text IS NULL OR transaction_type::text = $2)
	AND ($3::date IS NULL OR created_at::date >= $3)
	AND ($4::date IS NULL OR created_at::date <= $4)
	LIMIT $5 OFFSET $6`,
		param.AccountID, typeFilter, param.StartDate, param.EndDate, param.Limit, offset,
	)
	if err != nil {
		log.Printf("failed to get transactions for account %d: %v", param.AccountID, err)
		return nil, 0, appError.ErrInternalServer
	}
	var transactions []*models.Transaction
	defer rows.Close()
	for rows.Next() {
		var t models.Transaction
		var typeStr string
		if err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.ToAccountID,
			&t.Amount,
			&typeStr,
			&t.CreatedAt,
		); err != nil {
			log.Printf("failed to scan transaction for account %d: %v", param.AccountID, err)
			return nil, 0, appError.ErrInternalServer
		}
		t.Type = models.TransactionType(typeStr)
		transactions = append(transactions, &t)
	}
	if err := rows.Err(); err != nil {
		log.Printf("failed to get transactions for account %d: %v", param.AccountID, err)
		return nil, 0, appError.ErrInternalServer
	}
	return transactions, total, nil
}

func (r *TransactionRepository) lockAccount(ctx context.Context, tx pgx.Tx, id int64) error {
	_, err := tx.Exec(ctx, `UPDATE accounts SET is_locked = true WHERE id = $1`, id)
	return err
}

func (r *TransactionRepository) unlockAccount(ctx context.Context, tx pgx.Tx, id int64) error {
	_, err := tx.Exec(ctx, `UPDATE accounts SET is_locked = false WHERE id = $1`, id)
	return err
}
