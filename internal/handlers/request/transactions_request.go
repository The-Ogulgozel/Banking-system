package request

import (
	"github.com/The-Ogulgozel/Banking-system/internal/errors"
	"github.com/The-Ogulgozel/Banking-system/internal/validations"
)

type TransactionRequest struct {
	Amount int64 `json:"amount"`
}

func (req *TransactionRequest) Validate() error {
	if err := validations.ValidateAmount(req.Amount); err != nil {
		return err
	}
	return nil
}

type TransferRequest struct {
	Amount      int64 `json:"amount"`
	ToAccountID int64 `json:"to_account_id"`
}

func (req *TransferRequest) Validate() error {
	if err := validations.ValidateAmount(req.Amount); err != nil {
		return err
	}
	if req.ToAccountID <= 0 {
		return errors.ErrNotAccountId
	}
	return nil
}

type ListTransactions struct {
	Type      *string `form:"type"`
	Page      *int64  `form:"page"`
	Limit     *int64  `form:"limit"`
	StartDate *string `form:"start_date"`
	EndDate   *string `form:"end_date"`
}

func (req *ListTransactions) Validate() error {
	if req.Type != nil {
		if err := validations.ValidatetransactionType(*req.Type); err != nil {
			return err
		}
	}
	if err := validations.ValidateDateRange(req.StartDate, req.EndDate); err != nil {
		return err
	}
	if err := validations.ValidatePagination(req.Page, req.Limit); err != nil {
		return err
	}
	return nil
}
