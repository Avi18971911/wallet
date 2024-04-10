package services

import "context"

type TransactionServiceImpl struct {
}

func CreateNewTransactionServiceImpl() *TransactionServiceImpl {
	return &TransactionServiceImpl{}
}

func (t *TransactionServiceImpl) AddTransaction(
	toAccount string,
	fromAccount string,
	amount float64,
	ctx context.Context,
) {
	// do nothing
}
