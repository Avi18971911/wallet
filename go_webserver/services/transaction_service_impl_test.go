package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"webserver/mocks"
)

func TestAddTransaction(t *testing.T) {

	t.Run("Doesn't return error if transaction is successful", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockTran.On("Commit", mock.Anything).Return(nil)

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		assert.Nil(t, err)
	})

	t.Run("Returns error if repository Add Transaction is unsuccessful", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		assert.NotNil(t, err)
	})

	t.Run("Returns error if repository Add Balance is unsuccessful", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		assert.NotNil(t, err)
	})

	t.Run("Returns error if repository Deduct Balance is unsuccessful", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		assert.NotNil(t, err)
	})

	t.Run("Commits if no errors are encountered", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockTran.On("Commit", mock.Anything).Return(nil)

		service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		mockTran.AssertNumberOfCalls(t, "Commit", 1)
		mockTran.AssertNumberOfCalls(t, "Rollback", 0)
	})

	t.Run("Rollback if errors are encountered", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		mockTran.AssertNumberOfCalls(t, "Commit", 0)
		mockTran.AssertNumberOfCalls(t, "Rollback", 1)
	})

	t.Run("Returns error if error is encountered during commit", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockTran.On("Commit", mock.Anything).Return(errors.New("commit error"))

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "commit error")
		mockTran.AssertNumberOfCalls(t, "Commit", 1)
		mockTran.AssertNumberOfCalls(t, "Rollback", 0)
	})
}

func initializeMocks() (
	*mocks.MockTransactionRepository,
	*mocks.MockAccountRepository,
	*mocks.MockTransactional,
	*TransactionServiceImpl,
	context.Context,
	context.CancelFunc,
) {
	mockTranRepo := new(mocks.MockTransactionRepository)
	mockAccRepo := &mocks.MockAccountRepository{}
	mockTran := &mocks.MockTransactional{}

	service := CreateNewTransactionServiceImpl(mockTranRepo, mockAccRepo, mockTran)
	addCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return mockTranRepo, mockAccRepo, mockTran, service, addCtx, cancel
}
