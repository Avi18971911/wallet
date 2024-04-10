package services

import (
	"context"
	"webserver/domain"
	"webserver/repositories"
)

type TransactionServiceImpl struct {
	tr repositories.TransactionRepository
}

func CreateNewTransactionServiceImpl(tr repositories.TransactionRepository) *TransactionServiceImpl {
	return &TransactionServiceImpl{tr}
}

func (t *TransactionServiceImpl) AddTransaction(
	toAccount string,
	fromAccount string,
	amount float64,
	ctx context.Context,
) {
	t.tr.AddTransaction(domain.TransactionDetails{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
	}, ctx)
}
