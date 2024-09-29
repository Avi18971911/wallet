package model

import "github.com/shopspring/decimal"

type TransactionDetails struct {
	FromBankAccountId string
	ToBankAccountId   string
	Amount            decimal.Decimal
}
