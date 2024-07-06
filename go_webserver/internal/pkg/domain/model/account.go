package model

import "time"

type AccountDetails struct {
	Id               string
	Username         string
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
