package repositories

import (
	"context"
	domain2 "webserver/internal/pkg/domain"
)

type TransactionRepository interface {
	AddTransaction(details domain2.TransactionDetails, ctx context.Context) error
	GetAccountTransactions(accountId string, ctx context.Context) ([]domain2.AccountTransaction, error)
}
