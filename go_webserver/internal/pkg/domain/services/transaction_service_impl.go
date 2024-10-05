package services

import (
	"context"
	"fmt"
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

func (t *TransactionServiceImpl) AddTransaction(input model.TransactionDetailsInput, ctx context.Context) error {
	addCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := t.tran.BeginTransaction(addCtx, transactional.IsolationLow, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Add Transaction database transaction from "+
			"BankAccount %s to BankAccount %s: %v", input.FromBankAccountId, input.ToBankAccountId, err)
		return fmt.Errorf("error when starting Add Transaction database transaction: %w", err)
	}

	var shouldRollback = true

	defer func() {
		if shouldRollback {
			if rollErr := t.tran.Rollback(txnCtx); rollErr != nil {
				log.Printf("Error rolling back transaction: %v", rollErr)
			}
			return
		}
	}()

	toPending := input.Type == model.Pending
	newBalance, pendingBalance, err := t.ar.DeductBalance(input.FromBankAccountId, input.Amount, toPending, txnCtx)
	if err != nil {
		log.Printf("Error deducting balance from BankAccount %s: %v", input.FromBankAccountId, err)
		return fmt.Errorf("error when deducting balance from BankAccount %s: %w", input.FromBankAccountId, err)
	}

	if newBalance.IsNegative() || pendingBalance.IsNegative() {
		log.Printf("Insufficient balance in BankAccount %s", input.FromBankAccountId)
		return fmt.Errorf("insufficient balance in BankAccount %s", input.FromBankAccountId)
	}

	log.Printf("Successfully deducted balance from BankAccount %s", input.FromBankAccountId)

	if err = t.ar.AddBalance(input.ToBankAccountId, input.Amount, toPending, txnCtx); err != nil {
		log.Printf("Error adding balance to BankAccount %s: %v", input.ToBankAccountId, err)
		return fmt.Errorf("error when adding balance to BankAccount %s: %w", input.ToBankAccountId, err)
	}

	log.Printf("Successfully added balance to BankAccount %s", input.ToBankAccountId)

	if err = t.tr.AddTransaction(&input, txnCtx); err != nil {
		log.Printf("Error adding transaction to the database from BankAccount %s to "+
			"BankAccount %s: %v", input.FromBankAccountId, input.ToBankAccountId, err)
		return fmt.Errorf("error when adding transaction to the database: %w", err)
	}

	if commitErr := t.tran.Commit(txnCtx); commitErr != nil {
		log.Printf("Error committing Add Transaction database transaction: %v", commitErr)
		return fmt.Errorf("error when committing Add Transaction database transaction: %w", commitErr)
	}

	log.Printf(
		"Successfully committed Add Transaction database transaction from "+
			"BankAccount %s to BankAccount %s", input.FromBankAccountId, input.ToBankAccountId,
	)

	shouldRollback = false
	return nil
}
