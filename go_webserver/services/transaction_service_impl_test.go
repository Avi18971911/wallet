package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"webserver/mocks"
)

func TestAddTransaction(t *testing.T) {
	mockTranRepo := &mocks.MockTransactionRepository{}
	mockAccRepo := &mocks.MockAccountRepository{}
	mockTran := &mocks.MockTransactional{}

	service := CreateNewTransactionServiceImpl(mockTranRepo, mockAccRepo, mockTran)
	ctx := context.Background()
	addCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	t.Run("Doesn't return error if transaction is successful", func(t *testing.T) {
		mockTran.On("BeginTransaction", mock.Anything).Return(addCtx, nil).Once()
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil).Once()
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()
		mockTran.On("Commit", mock.Anything).Return(nil).Once()

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, addCtx)
		assert.Nil(t, err)
	})

	t.Run("Returns error if transaction isn't successful", func(t *testing.T) {
		mockTran.On("BeginTransaction", mock.Anything).Return(addCtx, nil).Once()
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(assert.AnError).Once()
		mockTran.On("Rollback", mock.Anything).Return(nil).Once()

		err := service.AddTransaction("toAccountID", "fromAccountID", 100.00, addCtx)
		assert.NotNil(t, err)
	})
}
