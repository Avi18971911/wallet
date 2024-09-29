package services

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type AccountService interface {
	GetAccountDetailsFromBankAccountId(bankAccountId string, ctx context.Context) (*model.AccountDetails, error)
	GetBankAccountTransactions(bankAccountId string, ctx context.Context) ([]model.BankAccountTransaction, error)
	Login(username string, password string, ctx context.Context) (*model.AccountDetails, error)
}
