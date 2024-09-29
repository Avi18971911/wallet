package handlers

import (
	"github.com/shopspring/decimal"
	"log"
	"webserver/internal/app/server/dto"
	"webserver/internal/pkg/domain/model"
)

func transactionDetailsToModel(tx *dto.TransactionRequest) (model.TransactionDetailsInput, error) {
	decimalAmount, err := decimal.NewFromString(tx.Amount)
	if err != nil {
		log.Printf("Failed to convert amount %s to decimal: %v", tx.Amount, err)
		return model.TransactionDetailsInput{}, err
	}

	return model.TransactionDetailsInput{
		FromBankAccountId: tx.FromBankAccountId,
		ToBankAccountId:   tx.ToBankAccountId,
		Amount:            decimalAmount,
		Type:              model.Realized,
	}, nil
}
