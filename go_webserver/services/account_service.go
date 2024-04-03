package services

import (
	"webserver/domain"
)

type AccountService interface {
	GetAccountDetails(accountId string) *domain.AccountDetails
	GetAccountTransactions(accountId string) []*domain.AccountTransaction
}
