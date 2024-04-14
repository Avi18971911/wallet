package services

import (
	"context"
	"webserver/internal/pkg/domain"
)

type AccountService interface {
	GetAccountDetails(accountId string, ctx context.Context) (*domain.AccountDetails, error)
	GetAccountTransactions(accountId string, ctx context.Context) ([]domain.AccountTransaction, error)
}
