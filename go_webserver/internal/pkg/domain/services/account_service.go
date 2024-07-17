package services

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type AccountService interface {
	GetAccountDetails(accountId string, ctx context.Context) (*model.AccountDetails, error)
	GetAccountTransactions(accountId string, ctx context.Context) ([]model.AccountTransaction, error)
	Login(username string, password string, ctx context.Context) (*model.AccountDetails, error)
}
