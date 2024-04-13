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
	return []domain.AccountTransaction{}, args.Error(0)
}
