package services

import (
	"context"
	"log"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/infrastructure/transactional"
)

type AccountServiceImpl struct {
	ar   repositories.AccountRepository
	tr   repositories.TransactionRepository
	tran transactional.Transactional
}

func CreateNewAccountServiceImpl(
	ar repositories.AccountRepository,
	tr repositories.TransactionRepository,
	tran transactional.Transactional,
) *AccountServiceImpl {
	return &AccountServiceImpl{ar: ar, tr: tr, tran: tran}
}

func (a *AccountServiceImpl) GetAccountDetails(accountId string, ctx context.Context) (*model.AccountDetails, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()
	accountDetails, err := a.ar.GetAccountDetails(accountId, getCtx)
	if err != nil {
		log.Printf("Unable to get account details for Account %s with error: %v", accountId, err)
		return nil, err
	}
	return accountDetails, nil
}

func (a *AccountServiceImpl) GetAccountTransactions(
	accountId string, ctx context.Context,
) ([]model.AccountTransaction, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := a.tran.BeginTransaction(getCtx, transactional.IsolationHigh, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Get Account Transactions database transaction for "+
			"Account %s: %v", accountId, err)
		return nil, err
	}

	defer func() {
		if rollErr := a.tran.Rollback(txnCtx); rollErr != nil {
			log.Printf("Error rolling back transaction: %v", rollErr)
		}
	}()

	accountTransactions, err := a.tr.GetAccountTransactions(accountId, getCtx)
	if err != nil {
		log.Printf("Unable to get transaction details for Account %s with error: %v", accountId, err)
		return nil, err
	}
	return accountTransactions, nil
}
