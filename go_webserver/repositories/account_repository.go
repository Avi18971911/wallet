package repositories

import (
	"context"
	"webserver/domain"
)

type AccountRepository interface {
	GetAccountDetails(accountId string, ctx context.Context) *domain.AccountDetails
	GetAccountTransactions(accountId string, ctx context.Context) []*domain.AccountTransaction
}
