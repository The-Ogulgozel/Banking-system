package models

import "time"

type Account struct {
	ID        int64
	Balance   int64
	Currency  string
	IsLocked  bool
	CreatedAt time.Time
}

type ListAccounts struct {
	Currency *string
	Page     *int64
	Limit    *int64
}
