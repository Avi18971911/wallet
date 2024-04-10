package services

import (
	"context"
	"webserver/domain"
	"webserver/repositories"
)

type AccountServiceImpl struct {
	repo repositories.AccountRepository
}

func CreateNewAccountServiceImpl(repo repositories.AccountRepository) *AccountServiceImpl {
	return &AccountServiceImpl{repo: repo}
}

func (a *AccountServiceImpl) GetAccountDetails(accountId string, ctx context.Context) *domain.AccountDetails {
	return a.repo.GetAccountDetails(accountId)
}

func (a *AccountServiceImpl) GetAccountTransactions(
	accountId string, ctx context.Context,
) []*domain.AccountTransaction {
	return a.repo.GetAccountTransactions(accountId)
}
