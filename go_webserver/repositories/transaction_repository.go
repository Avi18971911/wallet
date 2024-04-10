package repositories

import (
	"context"
	"webserver/domain"
)

type TransactionRepository interface {
	AddTransaction(
		details domain.TransactionDetails,
		ctx context.Context,
	)
}
