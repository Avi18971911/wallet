package model

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type AccountDetails struct {
	Id            string
	Username      string
	Password      string
	Person        Person
	Accounts      []Account
	KnownAccounts []KnownAccount
	CreatedAt     time.Time
}

const (
	Savings int = iota
	Checking
	Investment
)

type Person struct {
	FirstName string
	LastName  string
}

type Account struct {
	Id               string
	AccountNumber    string
	AccountType      int
	AvailableBalance decimal.Decimal
}

type KnownAccount struct {
	Id            string
	AccountNumber string
	AccountHolder string
	AccountType   int
}

type AccountTransaction struct {
	Id              string
	AccountId       string
	OtherAccountId  string
	TransactionType string
	Amount          decimal.Decimal
	CreatedAt       time.Time
}

var (
	ErrNoMatchingUsername = errors.New("no matching username found for account")
	ErrInvalidCredentials = errors.New("invalid username or password")
)
