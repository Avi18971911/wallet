package repositories

import (
	"context"
	"github.com/shopspring/decimal"
	"webserver/internal/pkg/domain/model"
)

type AccountRepository interface {
	GetAccountDetails(accountId string, ctx context.Context) (*model.AccountDetails, error)
	AddBalance(accountId string, amount decimal.Decimal, ctx context.Context) error
	DeductBalance(accountId string, amount decimal.Decimal, ctx context.Context) error
	GetAccountDetailsFromUsername(username string, ctx context.Context) (*model.AccountDetails, error)
}
