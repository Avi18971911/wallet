package repositories

import (
	"context"
	"webserver/domain"
)

type AccountRepository interface {
	GetAccountDetails(accountId string, ctx context.Context) *domain.AccountDetails
	GetAccountTransactions(accountId string, ctx context.Context) []*domain.AccountTransaction
	AddBalance(accountId string, amount float64, ctx context.Context) error
	DeductBalance(accountId string, amount float64, ctx context.Context) error
}
