package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"webserver/internal/pkg/domain/model"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) AddTransaction(details *model.TransactionDetails, ctx context.Context) error {
	args := m.Called(details, ctx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAccountTransactions(
	accountId string,
	ctx context.Context,
) ([]model.AccountTransaction, error) {
	args := m.Called(accountId, ctx)
	var accountTransactions []model.AccountTransaction
	if args.Get(0) != nil {
		accountTransactions = args.Get(0).([]model.AccountTransaction)
	}
	return accountTransactions, args.Error(1)
}
