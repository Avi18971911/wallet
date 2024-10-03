package services

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/test/mocks"
)

func TestAddTransaction(t *testing.T) {
	testAmt, _ := decimal.NewFromString("100.00")
	input := model.TransactionDetailsInput{
		ToBankAccountId:   "toAccountID",
		FromBankAccountId: "fromAccountID",
		Amount:            testAmt,
	}

	t.Run("Doesn't return error if transaction is successful", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockTran.On("Commit", mock.Anything).Return(nil)

		err := service.AddTransaction(input, ctx)
		assert.Nil(t, err)
	})

	t.Run("Returns error immediately if starting transaction is unsuccessful", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).
			Return(ctx, errors.New("can't start transaction"))

		err := service.AddTransaction(input, ctx)
		assert.ErrorContains(t, err, "can't start transaction")
		mockTran.AssertNumberOfCalls(t, "Rollback", 0)
		mockTranRepo.AssertNumberOfCalls(t, "AddTransaction", 0)
	})

	t.Run("Returns error if repository Add Transaction is unsuccessful", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		err := service.AddTransaction(input, ctx)
		assert.NotNil(t, err)
	})

	t.Run("Returns error if repository Add Balance is unsuccessful", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		err := service.AddTransaction(input, ctx)
		assert.NotNil(t, err)
	})

	t.Run("Returns error if repository Deduct Balance is unsuccessful", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		err := service.AddTransaction(input, ctx)
		assert.NotNil(t, err)
	})

	t.Run("Commits if no errors are encountered", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockTran.On("Commit", mock.Anything).Return(nil)

		service.AddTransaction(input, ctx)
		mockTran.AssertNumberOfCalls(t, "Commit", 1)
		mockTran.AssertNumberOfCalls(t, "Rollback", 0)
	})

	t.Run("Rollback if errors are encountered", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(assert.AnError)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		service.AddTransaction(input, ctx)
		mockTran.AssertNumberOfCalls(t, "Commit", 0)
		mockTran.AssertNumberOfCalls(t, "Rollback", 1)
	})

	t.Run("Returns error if error is encountered during commit", func(t *testing.T) {
		mockTranRepo, mockAccRepo, mockTran, service, ctx, cancel := initializeTransactionMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("AddTransaction", mock.Anything, mock.Anything).Return(nil)
		mockAccRepo.On("AddBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockAccRepo.On("DeductBalance", mock.Anything, mock.Anything, mock.Anything).
			Return(nil)
		mockTran.On("Commit", mock.Anything).Return(errors.New("commit error"))

		err := service.AddTransaction(input, ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "commit error")
		mockTran.AssertNumberOfCalls(t, "Commit", 1)
		mockTran.AssertNumberOfCalls(t, "Rollback", 0)
	})
}

func initializeTransactionMocks() (
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
