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
	"webserver/test/utils"
)

func TestGetAccountDetails(t *testing.T) {
	t.Run("Doesn't return error if GetAccountDetailsFromBankAccountId is successful", func(t *testing.T) {
		_, mockAccRepo, _, service, ctx, cancel := initializeAccountMocks()
		defer cancel()
		stubDetails := &utils.TomAccountDetailsModel
		mockAccRepo.On("GetAccountDetailsFromBankAccountId", mock.Anything, mock.Anything).Return(stubDetails, nil)

		res, err := service.GetAccountDetailsFromBankAccountId("accountId", ctx)
		assert.Nil(t, err)
		assert.EqualExportedValues(t, stubDetails, res, "The returned account details should match "+
			"the expected stub details")
	})

	t.Run("Returns error if GetAccountDetailsFromBankAccountId is unsuccessful", func(t *testing.T) {
		_, mockAccRepo, _, service, ctx, cancel := initializeAccountMocks()
		defer cancel()
		mockAccRepo.On("GetAccountDetailsFromBankAccountId", mock.Anything, mock.Anything).
			Return(nil, errors.New("cannot GetAccountDetailsFromBankAccountId"))

		res, err := service.GetAccountDetailsFromBankAccountId("accountId", ctx)
		assert.Nil(t, res)
		assert.Error(t, err, "cannot GetAccountDetailsFromBankAccountId")
	})
}

func TestGetAccountTransaction(t *testing.T) {
	t.Run("Returns correct output assuming happy path", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeAccountMocks()
		defer cancel()
		stubTransactions := []model.BankAccountTransactionOutput{
			{Id: "transactionId", BankAccountId: "accountId", Amount: decimal.NewFromFloat(123.13), CreatedAt: time.Now()},
		}
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTranRepo.On("GetBankAccountTransactions", mock.Anything, mock.Anything).
			Return(stubTransactions, nil)
		mockTran.On("Rollback", mock.Anything).Return(nil)

		accountTransactions, err := service.GetBankAccountTransactions("accountId", ctx)
		assert.Nil(t, err)
		assert.Equal(t, stubTransactions, accountTransactions)
	})

	t.Run("Returns error immediately if starting transaction is unsuccessful", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeAccountMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).
			Return(ctx, errors.New("can't start transaction"))

		_, err := service.GetBankAccountTransactions("accountId", ctx)
		assert.ErrorContains(t, err, "can't start transaction")
		mockTran.AssertNumberOfCalls(t, "Rollback", 0)
		mockTranRepo.AssertNumberOfCalls(t, "GetBankAccountTransactions", 0)
	})

	t.Run("Returns error if GetBankAccountTransactions isn't successful", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeAccountMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTran.On("Rollback", mock.Anything).Return(nil)
		mockTranRepo.On("GetBankAccountTransactions", mock.Anything, mock.Anything).
			Return(nil, errors.New("can't GetBankAccountTransactions"))

		_, err := service.GetBankAccountTransactions("accountId", ctx)
		assert.ErrorContains(t, err, "can't GetBankAccountTransactions")
		mockTranRepo.AssertNumberOfCalls(t, "GetBankAccountTransactions", 1)
	})

	t.Run("Rollback even if returning an error", func(t *testing.T) {
		mockTranRepo, _, mockTran, service, ctx, cancel := initializeAccountMocks()
		defer cancel()
		mockTran.On("BeginTransaction", mock.Anything, mock.Anything, mock.Anything).Return(ctx, nil)
		mockTran.On("Rollback", mock.Anything).Return(nil)
		mockTranRepo.On("GetBankAccountTransactions", mock.Anything, mock.Anything).
			Return(nil, errors.New("can't GetBankAccountTransactions"))

		service.GetBankAccountTransactions("accountId", ctx)
		mockTran.AssertNumberOfCalls(t, "Rollback", 1)
	})
}

func initializeAccountMocks() (
	*mocks.MockTransactionRepository,
	*mocks.MockAccountRepository,
	*mocks.MockTransactional,
	*AccountServiceImpl,
	context.Context,
	context.CancelFunc,
) {
	mockTranRepo := new(mocks.MockTransactionRepository)
	mockAccRepo := &mocks.MockAccountRepository{}
	mockTran := &mocks.MockTransactional{}

	service := CreateNewAccountServiceImpl(mockAccRepo, mockTranRepo, mockTran)
	addCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	return mockTranRepo, mockAccRepo, mockTran, service, addCtx, cancel
}
