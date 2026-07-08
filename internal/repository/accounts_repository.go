package repository

import (
	"context"
	"log"

	"errors"

	appError "github.com/The-Ogulgozel/Banking-system/internal/errors"

	"github.com/The-Ogulgozel/Banking-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepo(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Create(ctx context.Context, dto *models.Account) error {
	query := `
	INSERT INTO accounts (balance, currency)
	VALUES($1,$2)
	RETURNING id,is_locked,created_at
	`
	err := r.db.QueryRow(ctx, query,
		dto.Balance,
		dto.Currency,
	).Scan(&dto.ID, &dto.IsLocked, &dto.CreatedAt)

	if err != nil {
		log.Printf("failed to create account: %v", err)
		return appError.ErrInternalServer
	}
	return nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int64) (*models.Account, error) {
	query := `
	SELECT id,balance,currency,is_locked,created_at
	FROM accounts 
	WHERE id=$1 AND deleted_at IS NULL`

	var a models.Account
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID,
		&a.Balance,
		&a.Currency,
		&a.IsLocked,
		&a.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appError.ErrNotFound
		}
		log.Printf("failed to get account by id: %v", err)
		return nil, appError.ErrInternalServer
	}
	return &a, nil
}

func (r *AccountRepository) GetAll(ctx context.Context, param *models.ListAccounts) ([]*models.Account, int64, error) {

	offset := (*param.Page - 1) * *param.Limit

	var total int64
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		 FROM accounts 
		 WHERE  deleted_at IS NULL AND
		 ($1::text IS NULL OR currency = $1)`,
		param.Currency,
	).Scan(&total)

	if err != nil {
		log.Printf("failed to get total accounts: %v", err)
		return nil, 0, appError.ErrInternalServer
	}

	query := `
	SELECT id,balance,currency,is_locked,created_at 
	FROM accounts 
	WHERE deleted_at IS NULL AND 
	($1::text IS NULL OR currency = $1) 
	ORDER BY id LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, param.Currency, *param.Limit, offset)
	if err != nil {
		log.Printf("failed to get accounts: %v", err)
		return nil, 0, appError.ErrInternalServer
	}
	var accounts []*models.Account
	defer rows.Close()
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(
			&a.ID,
			&a.Balance,
			&a.Currency,
			&a.IsLocked,
			&a.CreatedAt,
		); err != nil {
			log.Printf("failed to scan accounts: %v", err)
			return nil, 0, appError.ErrInternalServer
		}
		accounts = append(accounts, &a)
	}
	if err := rows.Err(); err != nil {
		log.Printf("failed to scan accounts: %v", err)
		return nil, 0, appError.ErrInternalServer
	}
	return accounts, total, nil
}

func (r *AccountRepository) Delete(ctx context.Context, id int64) error {
	query := `UPDATE accounts SET deleted_at=now() WHERE id=$1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		log.Printf("failed to delete account: %v", err)
		return appError.ErrInternalServer
	}
	if result.RowsAffected() == 0 {
		return appError.ErrNotFound
	}
	return nil
}
