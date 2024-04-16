package repositories

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type AccountRepository interface {
	GetAccountDetails(accountId string, ctx context.Context) (*model.AccountDetails, error)
	AddBalance(accountId string, amount float64, ctx context.Context) error
	DeductBalance(accountId string, amount float64, ctx context.Context) error
}
