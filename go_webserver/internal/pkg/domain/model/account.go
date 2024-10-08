package model

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type AccountDetailsOutput struct {
	Id                string
	Username          string
	Password          string
	Person            Person
	BankAccounts      []BankAccount
	KnownBankAccounts []KnownBankAccount
	CreatedAt         time.Time
}

type BankAccountType string

const (
	Savings    BankAccountType = "savings"
	Checking   BankAccountType = "checking"
	Investment BankAccountType = "investment"
)

type Person struct {
	FirstName string
	LastName  string
}

type BankAccount struct {
	Id               string
	AccountNumber    string
	AccountType      BankAccountType
	PendingBalance   decimal.Decimal
	AvailableBalance decimal.Decimal
}

type KnownBankAccount struct {
	Id            string
	AccountNumber string
	AccountHolder string
	AccountType   BankAccountType
}

type BankAccountTransactionOutput struct {
	Id                 string
	BankAccountId      string
	OtherBankAccountId string
	TransactionNature  TransactionNature
	TransactionType    TransactionType
	ExpirationDate     time.Time
	Status             PendingTransactionStatus
	Amount             decimal.Decimal
	CreatedAt          time.Time
}

var (
	ErrNoMatchingUsername = errors.New("no matching username found for account")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type AccountBalanceMonthsOutput struct {
	BankAccountId string
	Months        []AccountBalanceMonth
}

type AccountBalanceMonth struct {
	Month            time.Month
	Year             int
	AvailableBalance decimal.Decimal
	PendingBalance   decimal.Decimal
}
