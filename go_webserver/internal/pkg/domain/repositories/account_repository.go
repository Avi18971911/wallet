package repositories

import (
	"context"
	"github.com/shopspring/decimal"
	"webserver/internal/pkg/domain/model"
)

type AccountRepository interface {
	GetAccountDetailsFromBankAccountId(bankAccountId string, ctx context.Context) (*model.AccountDetailsOutput, error)
	AddBalance(bankAccountId string, amount decimal.Decimal, toPending bool, ctx context.Context) error
	DeductBalance(bankAccountId string, amount decimal.Decimal, toPending bool, ctx context.Context) (
		decimal.Decimal,
		decimal.Decimal,
		error,
	)
	GetAccountDetailsFromUsername(username string, ctx context.Context) (*model.AccountDetailsOutput, error)
	GetAccountBalance(bankAccountId string, ctx context.Context) (decimal.Decimal, decimal.Decimal, error)
}
