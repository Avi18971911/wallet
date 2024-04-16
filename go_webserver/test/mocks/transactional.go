package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"webserver/internal/pkg/infrastructure/transactional"
)

type MockTransactional struct {
	mock.Mock
}

func (m *MockTransactional) BeginTransaction(ctx context.Context) (transactional.TransactionContext, error) {
	args := m.Called(ctx)
	return args.Get(0).(context.Context), args.Error(1)
}

func (m *MockTransactional) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTransactional) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
