package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"webserver/domain"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) AddTransaction(details domain.TransactionDetails, ctx context.Context) error {
	args := m.Called(details, ctx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAccountTransactions(
	accountId string,
	ctx context.Context,
) ([]domain.AccountTransaction, error) {
	args := m.Called(accountId, ctx)
	var accountTransactions []domain.AccountTransaction
	if args.Get(0) != nil {
		accountTransactions = args.Get(0).([]domain.AccountTransaction)
	}
	return accountTransactions, args.Error(1)
}
