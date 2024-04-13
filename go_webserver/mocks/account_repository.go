package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"webserver/domain"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccountDetails(
	accountId string,
	ctx context.Context,
) (*domain.AccountDetails, error) {
	args := m.Called(accountId, ctx)
	return &domain.AccountDetails{}, args.Error(0)
}

func (m *MockAccountRepository) AddBalance(accountID string, amount float64, ctx context.Context) error {
	args := m.Called(accountID, amount, ctx)
	return args.Error(0)
}

func (m *MockAccountRepository) DeductBalance(accountID string, amount float64, ctx context.Context) error {
	args := m.Called(accountID, amount, ctx)
	return args.Error(0)
}
