package services

import (
	"context"
	"log"
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
	accountDetails, err := a.repo.GetAccountDetails(accountId, ctx)
	if err != nil {
		log.Printf("Unable to get account details for Account %s with error: %v", accountId, err)
	}
	return accountDetails
}

func (a *AccountServiceImpl) GetAccountTransactions(
	accountId string, ctx context.Context,
) []*domain.AccountTransaction {
	accountTransactions, err := a.repo.GetAccountTransactions(accountId, ctx)
	if err != nil {
		log.Printf("Unable to get transaction details for Account %s with error: %v", accountId, err)
	}
	return accountTransactions
}
