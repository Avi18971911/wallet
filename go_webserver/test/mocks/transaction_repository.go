package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	domain2 "webserver/internal/pkg/domain"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) AddTransaction(details domain2.TransactionDetails, ctx context.Context) error {
	args := m.Called(details, ctx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetAccountTransactions(
	accountId string,
	ctx context.Context,
) ([]domain2.AccountTransaction, error) {
	args := m.Called(accountId, ctx)
	var accountTransactions []domain2.AccountTransaction
	if args.Get(0) != nil {
		accountTransactions = args.Get(0).([]domain2.AccountTransaction)
	}
	return accountTransactions, args.Error(1)
}
