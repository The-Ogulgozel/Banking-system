package response

import "github.com/The-Ogulgozel/Banking-system/internal/models"

type TransactionResponse struct {
	ID            int64  `json:"id"`
	FromAccountID int64  `json:"account_id"`
	ToAccountID   *int64 `json:"to_account_id,omitempty"`
	Type          string `json:"type"`
	Amount        int64  `json:"amount"`
	CreatedAt     string `json:"created_at"`
}

func ToTransactionResponse(t *models.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:            t.ID,
		FromAccountID: t.AccountID,
		ToAccountID:   t.ToAccountID,
		Type:          string(t.Type),
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt.Format("2006-01-02"),
	}
}

type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Page         int64                 `json:"page"`
	Limit        int64                 `json:"limit"`
}

func ToListTransactionsResponse(transactions []*models.Transaction, total, page, limit int64) ListTransactionsResponse {

	response := ListTransactionsResponse{
		Transactions: make([]TransactionResponse, 0, len(transactions)),
	}
	for _, transaction := range transactions {
		response.Transactions = append(response.Transactions, ToTransactionResponse(transaction))
	}
	response.Total = total
	response.Page = page
	response.Limit = limit
	return response
}
