package response

import "github.com/The-Ogulgozel/Banking-system/internal/models"

type AccountResponse struct {
	ID        int64  `json:"id"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
	IsLocked  bool   `json:"is_locked"`
	CreatedAt string `json:"created_at"`
}

func ToAccounResponse(a *models.Account) AccountResponse {
	response := AccountResponse{
		ID:        a.ID,
		Balance:   a.Balance,
		Currency:  a.Currency,
		IsLocked:  a.IsLocked,
		CreatedAt: a.CreatedAt.Format("2006-06-06"),
	}
	return response
}

type ListAccountsResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int64             `json:"total"`
	Page     int64             `json:"page"`
	Limit    int64             `json:"limit"`
}

func ToListAccountsResponse(accounts []*models.Account, total, page, limit int64) ListAccountsResponse {
	response := ListAccountsResponse{
		Accounts: make([]AccountResponse, 0, len(accounts)),
	}
	for _, account := range accounts {
		response.Accounts = append(response.Accounts, ToAccounResponse(account))
	}
	response.Total = total
	response.Page = page
	response.Limit = limit
	return response
}
