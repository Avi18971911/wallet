package repositories

import "webserver/domain"

type AccountRepository interface {
	GetAccountDetails(accountId string) *domain.AccountDetails
	GetAccountTransactions(accountId string) []*domain.AccountTransaction
}
