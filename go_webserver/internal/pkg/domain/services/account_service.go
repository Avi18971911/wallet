package services

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type AccountService interface {
	GetAccountDetailsFromBankAccountId(bankAccountId string, ctx context.Context) (*model.AccountDetailsOutput, error)
	GetBankAccountTransactions(
		input *model.TransactionsForBankAccountInput,
		ctx context.Context,
	) (
		[]model.BankAccountTransactionOutput, error)
	Login(username string, password string, ctx context.Context) (*model.AccountDetailsOutput, error)
	GetAccountHistoryInMonths(
		input *model.AccountHistoryInMonthsInput,
		ctx context.Context,
	) (model.AccountBalanceMonthsOutput, error)
}
