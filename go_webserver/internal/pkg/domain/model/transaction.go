package model

import "github.com/shopspring/decimal"

type TransactionDetails struct {
	FromAccount string
	ToAccount   string
	Amount      decimal.Decimal
}
