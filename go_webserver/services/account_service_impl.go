package services

import (
	"webserver/domain"
	"webserver/repositories"
)

type AccountServiceImpl struct {
	repo repositories.AccountRepository
}

func CreateNewAccountServiceImpl(repo repositories.AccountRepository) *AccountServiceImpl {
	return &AccountServiceImpl{repo: repo}
}

func (a *AccountServiceImpl) GetAccountDetails(accountId string) *domain.AccountDetails {
	deets := a.repo.GetAccountDetails(accountId)
	return deets
}

func (a *AccountServiceImpl) GetAccountTransactions(accountId string) []*domain.AccountTransaction {
	transactions := a.repo.GetAccountTransactions(accountId)
	return transactions
}
