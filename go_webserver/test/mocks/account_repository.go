package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"webserver/internal/pkg/domain"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccountDetails(
	accountId string,
	ctx context.Context,
) (*domain.AccountDetails, error) {
	args := m.Called(accountId, ctx)
	var accountDetails *domain.AccountDetails
	if args.Get(0) != nil {
		accountDetails = args.Get(0).(*domain.AccountDetails)
	}
	return accountDetails, args.Error(1)
}

func (m *MockAccountRepository) AddBalance(accountID string, amount float64, ctx context.Context) error {
	args := m.Called(accountID, amount, ctx)
	return args.Error(0)
}

func (m *MockAccountRepository) DeductBalance(accountID string, amount float64, ctx context.Context) error {
	args := m.Called(accountID, amount, ctx)
	return args.Error(0)
}
