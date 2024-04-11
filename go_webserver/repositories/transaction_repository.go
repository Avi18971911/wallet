package repositories

import (
	"context"
	"webserver/domain"
)

type TransactionRepository interface {
	AddTransaction(details domain.TransactionDetails, ctx context.Context) error
	GetAccountTransactions(accountId string, ctx context.Context) ([]domain.AccountTransaction, error)
}
