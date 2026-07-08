package errors

type AppError struct {
	Err     error
	Message string
	Code    string
}

var (
	ErrNotFound               = NewAppError(nil, "not found", "SE-00404")
	ErrInternalServer         = NewAppError(nil, "internal server error", "SE-00500")
	ErrAccountLocked          = NewAppError(nil, "account id locked", "SE-00409")
	ErrAccountBalanceNotEmpty = NewAppError(nil, "account balance is not 0", "SE-00400")
	ErrInvalidCurrency        = NewAppError(nil, "currency is not supported", "SE-00400")
	ErrInvalidRequest         = NewAppError(nil, "invalid request", "SE-00400")
	ErrNotAccountId           = NewAppError(nil, "not valid account id", "SE-00400")
	ErrInvalidAmount          = NewAppError(nil, "amount must be positive and less than 1 trillion", "SE-00400")
	ErrNotEnoughBalance       = NewAppError(nil, "not enough balance", "SE-00400")
	ErrCurrencyMismatch       = NewAppError(nil, "currency mismatch", "SE-00400")
	ErrBalanceLimitExceeded   = NewAppError(nil, "balance limit exceeded", "SE-00400")
)

func NewAppError(err error, message, code string) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
func (e *AppError) Error() string {
	return e.Message
}
