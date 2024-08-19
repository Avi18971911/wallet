package model

import (
	"errors"
	"time"
)

type AccountDetails struct {
	Id               string
	Username         string
	Password         string
	AvailableBalance float64
	AccountNumber    string
	AccountType      int
	Person           Person
	KnownAccounts    []KnownAccount
	CreatedAt        time.Time
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

type KnownAccount struct {
	AccountNumber string
	AccountHolder string
	AccountType   int
}

type AccountTransaction struct {
	Id              string
	AccountId       string
	OtherAccountId  string
	TransactionType string
	Amount          float64
	CreatedAt       time.Time
}

var (
	ErrNoMatchingUsername = errors.New("no matching username found for account")
	ErrInvalidCredentials = errors.New("invalid username or password")
)
