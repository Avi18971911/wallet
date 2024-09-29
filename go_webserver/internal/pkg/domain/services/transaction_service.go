package services

import (
	"context"
)

type TransactionService interface {
	AddTransaction(
		toBankAccountId string,
		fromBankAccountId string,
		amount string,
		ctx context.Context,
	) error
}
