package repositories

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type TransactionRepository interface {
	AddTransaction(details *model.TransactionDetailsInput, ctx context.Context) error
	GetTransactionsFromBankAccountId(bankAccountId string, ctx context.Context) ([]model.BankAccountTransactionOutput, error)
}
