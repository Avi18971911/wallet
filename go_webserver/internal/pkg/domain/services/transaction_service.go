package services

import (
	"context"
)

type TransactionService interface {
	AddTransaction(
		toAccount string,
		fromAccount string,
		amount string,
		ctx context.Context,
	) error
}
