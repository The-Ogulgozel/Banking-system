package request

import "github.com/The-Ogulgozel/Banking-system/internal/validations"

type CreateAccountRequest struct {
	Balance  int64  `json:"balance" example:"1000"`
	Currency string `json:"currency" binding:"required" example:"USD"`
}

func (req *CreateAccountRequest) Validate() error {
	if err := validations.ValidateBalance(req.Balance); err != nil {
		return err
	}
	if err := validations.ValidateCurrency(req.Currency); err != nil {
		return err
	}
	return nil
}

type ListAccountRequest struct {
	Currency *string `form:"currency" binding:"omitempty" example:"USD"`
	Page     *int64  `form:"page" binding:"omitempty" example:"1"`
	Limit    *int64  `form:"limit" binding:"omitempty" example:"10"`
}

func (req *ListAccountRequest) Validate() error {
	if req.Currency != nil {
		if err := validations.ValidateCurrency(*req.Currency); err != nil {
			return err
		}
	}
	if err := validations.ValidatePagination(req.Page, req.Limit); err != nil {
		return err
	}
	return nil
}
