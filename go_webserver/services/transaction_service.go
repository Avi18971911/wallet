package services

import "context"

type TransactionService interface {
	AddTransaction(toAccount string, fromAccount string, amount float64, ctx context.Context)
}
