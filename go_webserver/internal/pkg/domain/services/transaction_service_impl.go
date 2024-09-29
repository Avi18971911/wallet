package services

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
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
	toBankAccountId string,
	fromBankAccountId string,
	amount string,
	ctx context.Context,
) error {
	addCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := t.tran.BeginTransaction(addCtx, transactional.IsolationLow, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Add Transaction database transaction from "+
			"BankAccount %s to BankAccount %s: %v", fromBankAccountId, toBankAccountId, err)
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

	amountDecimal, err := convertStringToDecimal(amount)
	if err != nil {
		log.Printf("Error converting string to decimal: %v", err)
		return fmt.Errorf("error when converting string to decimal: %w", err)
	}

	transactionDetails := model.TransactionDetails{
		FromBankAccountId: fromBankAccountId,
		ToBankAccountId:   toBankAccountId,
		Amount:            amountDecimal,
	}

	newBalance, err := t.ar.DeductBalance(fromBankAccountId, amountDecimal, txnCtx)
	if err != nil {
		log.Printf("Error deducting balance from BankAccount %s: %v", fromBankAccountId, err)
		return fmt.Errorf("error when deducting balance from BankAccount %s: %w", fromBankAccountId, err)
	}

	if newBalance.IsNegative() {
		log.Printf("Insufficient balance in BankAccount %s", fromBankAccountId)
		return fmt.Errorf("insufficient balance in BankAccount %s", fromBankAccountId)
	}

	if err = t.ar.AddBalance(toBankAccountId, amountDecimal, txnCtx); err != nil {
		log.Printf("Error adding balance to BankAccount %s: %v", toBankAccountId, err)
		return fmt.Errorf("error when adding balance to BankAccount %s: %w", toBankAccountId, err)
	}

	if err = t.tr.AddTransaction(&transactionDetails, txnCtx); err != nil {
		log.Printf("Error adding transaction to the database from BankAccount %s to "+
			"BankAccount %s: %v", fromBankAccountId, toBankAccountId, err)
		return fmt.Errorf("error when adding transaction to the database: %w", err)
	}

	if commitErr := t.tran.Commit(txnCtx); commitErr != nil {
		log.Printf("Error committing Add Transaction database transaction: %v", commitErr)
		return fmt.Errorf("error when committing Add Transaction database transaction: %w", commitErr)
	}

	shouldRollback = false
	return nil
}

func convertStringToDecimal(amount string) (decimal.Decimal, error) {
	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("error converting string to decimal: %w", err)
	}
	return amountDecimal, nil
}
