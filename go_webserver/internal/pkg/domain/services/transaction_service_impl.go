package services

import (
	"context"
	"log"
	"webserver/internal/pkg/domain/model"
	repositories2 "webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/infrastructure/transactional"
)

type TransactionServiceImpl struct {
	tr   repositories2.TransactionRepository
	ar   repositories2.AccountRepository
	tran transactional.Transactional
}

func CreateNewTransactionServiceImpl(
	tr repositories2.TransactionRepository,
	ar repositories2.AccountRepository,
	transactional transactional.Transactional,
) *TransactionServiceImpl {
	return &TransactionServiceImpl{tr, ar, transactional}
}

func (t *TransactionServiceImpl) AddTransaction(
	toAccount string,
	fromAccount string,
	amount float64,
	ctx context.Context,
) error {
	addCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := t.tran.BeginTransaction(addCtx, transactional.IsolationLow, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Add Transaction database transaction from "+
			"Account %s to Account %s: %v", fromAccount, toAccount, err)
		return err
	}

	defer func() {
		if err != nil {
			if rollErr := t.tran.Rollback(txnCtx); rollErr != nil {
				log.Printf("Error rolling back transaction: %v", rollErr)
			}
			return
		}
	}()

	transactionDetails := model.TransactionDetails{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
	}
	if err = t.tr.AddTransaction(&transactionDetails, txnCtx); err != nil {
		log.Printf("Error adding transaction to the database from Account %s to "+
			"Account %s: %v", fromAccount, toAccount, err)
		return err
	}

	if err = t.ar.AddBalance(toAccount, amount, txnCtx); err != nil {
		log.Printf("Error adding balance to Account %s: %v", toAccount, err)
		return err
	}

	if err = t.ar.DeductBalance(fromAccount, amount, txnCtx); err != nil {
		log.Printf("Error deducting balance from Account %s: %v", fromAccount, err)
		return err
	}

	if commitErr := t.tran.Commit(txnCtx); commitErr != nil {
		log.Printf("Error committing Add Transaction database transaction: %v", commitErr)
		return commitErr
	}

	return nil
}
