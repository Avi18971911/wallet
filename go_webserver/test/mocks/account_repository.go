package mocks

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"webserver/internal/pkg/domain/model"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccountDetailsFromBankAccountId(
	accountId string,
	ctx context.Context,
) (*model.AccountDetailsOutput, error) {
	args := m.Called(accountId, ctx)
	var accountDetails *model.AccountDetailsOutput
	if args.Get(0) != nil {
		accountDetails = args.Get(0).(*model.AccountDetailsOutput)
	}
	return accountDetails, args.Error(1)
}

func (m *MockAccountRepository) AddBalance(
	bankAccountID string,
	amount decimal.Decimal,
	toPending bool,
	ctx context.Context,
) error {
	args := m.Called(bankAccountID, amount, toPending, ctx)
	return args.Error(0)
}

func (m *MockAccountRepository) DeductBalance(
	accountID string,
	amount decimal.Decimal,
	toPending bool,
	ctx context.Context,
) (decimal.Decimal, decimal.Decimal, error) {
	args := m.Called(accountID, amount, toPending, ctx)
	return args.Get(0).(decimal.Decimal), args.Get(1).(decimal.Decimal), args.Error(2)
}

func (m *MockAccountRepository) GetAccountDetailsFromUsername(
	username string,
	ctx context.Context,
) (*model.AccountDetailsOutput, error) {
	args := m.Called(username, ctx)
	var accountDetails *model.AccountDetailsOutput
	if args.Get(0) != nil {
		accountDetails = args.Get(0).(*model.AccountDetailsOutput)
	}
	return accountDetails, args.Error(1)
}

func (m *MockAccountRepository) GetAccountBalance(bankAccountId string, ctx context.Context) (
	decimal.Decimal,
	decimal.Decimal, error,
) {
	args := m.Called(bankAccountId, ctx)
	var balance decimal.Decimal
	if args.Get(0) != nil {
		balance = args.Get(0).(decimal.Decimal)
	}
	var pendingBalance decimal.Decimal
	if args.Get(1) != nil {
		pendingBalance = args.Get(1).(decimal.Decimal)
	}
	return balance, pendingBalance, args.Error(1)
}
