package services

import (
	"context"
	"webserver/internal/pkg/domain/model"
)

type TransactionService interface {
	AddTransaction(
		input model.TransactionDetailsInput,
		ctx context.Context,
	) error
}
