package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type TransactionDetailsInput struct {
	FromBankAccountId string
	ToBankAccountId   string
	Amount            decimal.Decimal
	Type              TransactionType
	ExpirationDate    time.Time
	Status            PendingTransactionStatus
}

type TransactionDetailsOutput struct {
	Id                string
	FromBankAccountId string
	ToBankAccountId   string
	Amount            decimal.Decimal
	Type              TransactionType
	ExpirationDate    time.Time
	Status            PendingTransactionStatus
}

type TransactionType string

const (
	Realized TransactionType = "realized"
	Pending  TransactionType = "pending"
)

type PendingTransactionStatus string

const (
	Active  PendingTransactionStatus = "active"
	Applied PendingTransactionStatus = "applied"
	Revoked PendingTransactionStatus = "revoked"
)

type TransactionForBankAccountInput struct {
	BankAccountId string
	FromTime      time.Time
	ToTime        time.Time
}
