package services

import (
	"context"
	"webserver/domain"
)

type AccountService interface {
	GetAccountDetails(accountId string, ctx context.Context) *domain.AccountDetails
	GetAccountTransactions(accountId string, ctx context.Context) []*domain.AccountTransaction
}
