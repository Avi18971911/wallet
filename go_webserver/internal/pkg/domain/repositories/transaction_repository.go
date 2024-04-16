package repositories

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type TransactionRepository interface {
	AddTransaction(details *model.TransactionDetails, ctx context.Context) error
	GetAccountTransactions(accountId string, ctx context.Context) ([]model.AccountTransaction, error)
}
