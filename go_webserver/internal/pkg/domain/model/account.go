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
	CreatedAt        time.Time
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
