package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"webserver/internal/pkg/domain/model"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) AddTransaction(details *model.TransactionDetailsInput, ctx context.Context) error {
	args := m.Called(details, ctx)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetTransactionsFromBankAccountId(
	bankAccountId string,
	ctx context.Context,
) ([]model.BankAccountTransactionOutput, error) {
	args := m.Called(bankAccountId, ctx)
	var accountTransactions []model.BankAccountTransactionOutput
	if args.Get(0) != nil {
		accountTransactions = args.Get(0).([]model.BankAccountTransactionOutput)
	}
	return accountTransactions, args.Error(1)
}
