package validations

import (
	"time"

	"github.com/The-Ogulgozel/Banking-system/internal/errors"
)

func ValidateAmount(amount int64) error {
	if amount <= 0 {
		return errors.NewAppError(nil, "amount must be greater than zero", "SE-00400")
	}
	if amount > 1_000_000_000_000 {
		return errors.NewAppError(nil, "amount cannot be greater than 1 trillion", "SE-00400")
	}
	return nil
}

func ValidateDate(dateStr *string) error {
	if dateStr == nil {
		return nil
	}

	_, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return errors.NewAppError(nil, "invalid date format", "SE-00400")
	}

	return nil
}

func ValidateDateRange(beginDate, endDate *string) error {
	if err := ValidateDate(beginDate); err != nil {
		return err
	}

	if err := ValidateDate(endDate); err != nil {
		return err
	}
	if beginDate != nil && endDate != nil {
		if *beginDate > *endDate {
			return errors.NewAppError(nil, "begin date must be before or equal to end date", "SE-00400")
		}
	}
	return nil
}

func ValidatetransactionType(t string) error {
	if t != "deposit" && t != "withdraw" && t != "transfer" {
		return errors.NewAppError(nil, "invalid transaction type", "SE-00400")
	}
	return nil
}
