package validations

import "github.com/The-Ogulgozel/Banking-system/internal/errors"

var AllowedCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"JPY": true,
	"CNY": true,
	"CHF": true,
	"RUB": true,
	"TRY": true,
	"AZN": true,
	"TMT": true,
	"KZT": true,
	"UAH": true,
	"AED": true,
	"INR": true,
	"CAD": true,
}

func ValidateCurrency(currency string) error {
	if !AllowedCurrencies[currency] {
		return errors.ErrInvalidCurrency
	}
	return nil
}

func ValidateBalance(balance int64) error {
	if balance < 0 {
		return errors.NewAppError(nil, "balance cannot be negative", "SE-00400")
	}
	if balance > 1000_000_000_000 {
		return errors.NewAppError(nil, "balance cannot be greater than 1 trillion", "SE-00400")
	}
	return nil
}

func ValidatePagination(page, limit *int64) error {
	if page != nil {
		if *page < 1 {
			return errors.NewAppError(nil, "page must be at least 1", "SE-00400")
		}
		if *page > 10000 {
			return errors.NewAppError(nil, "page must be less than 10000", "SE-00400")
		}
	}
	if limit != nil {
		if *limit < 1 {
			return errors.NewAppError(nil, "limit must be at least 1", "SE-00400")
		}
		if *limit > 100 {
			return errors.NewAppError(nil, "limit must be less than or equal to 100", "SE-00400")
		}
	}
	return nil
}
