package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/infrastructure/transactional"
)

type AccountServiceImpl struct {
	ar   repositories.AccountRepository
	tr   repositories.TransactionRepository
	tran transactional.Transactional
}

func CreateNewAccountServiceImpl(
	ar repositories.AccountRepository,
	tr repositories.TransactionRepository,
	tran transactional.Transactional,
) *AccountServiceImpl {
	return &AccountServiceImpl{ar: ar, tr: tr, tran: tran}
}

func validateAccountDetails(accountDetails *model.AccountDetails) error {
	err := validateAccountNumbers(accountDetails.Accounts)
	if err != nil {
		return fmt.Errorf("unable to validate account numbers with error: %w", err)
	}
	err = validateAccountTypes(accountDetails.Accounts)
	if err != nil {
		return fmt.Errorf("unable to validate account types with error: %w", err)
	}
	return nil
}

func validateAccountNumbers(accounts []model.Account) error {
	const pattern = `^\d{3}-\d{5}-\d{1}$`
	var accountNumberRegex, err = regexp.Compile(pattern)
	if err != nil {
		log.Printf("Unable to compile regex pattern for account number validation with error: %v", err)
	} else {
		for _, account := range accounts {
			if !accountNumberRegex.MatchString(account.AccountNumber) {
				return errors.New("account number does not match expected pattern of XXX-XXXXX-X")
			}
		}
	}
	return nil
}

func validateAccountTypes(accounts []model.Account) error {
	for _, account := range accounts {
		if account.AccountType > model.Investment || account.AccountType < model.Savings {
			return errors.New("account type is not a valid type of account. " +
				"Should be savings, checking, or investment")
		}
	}
	return nil
}

func (a *AccountServiceImpl) GetAccountDetails(accountId string, ctx context.Context) (*model.AccountDetails, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()
	accountDetails, err := a.ar.GetAccountDetails(accountId, getCtx)
	if err != nil {
		log.Printf("Unable to get account details for Account %s with error: %v", accountId, err)
		return nil, fmt.Errorf("unable to get account details with error: %w", err)
	}
	err = validateAccountDetails(accountDetails)
	if err != nil {
		log.Printf("Unable to successfully validate account details for Account %s with error: "+
			"%v", accountId, err)
		return nil, fmt.Errorf("unable to validate account details with error: %w", err)
	}
	return accountDetails, nil
}

func (a *AccountServiceImpl) GetAccountTransactions(
	accountId string, ctx context.Context,
) ([]model.AccountTransaction, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := a.tran.BeginTransaction(getCtx, transactional.IsolationHigh, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Get Account Transactions database transaction for "+
			"Account %s: %v", accountId, err)
		return nil, fmt.Errorf("unable to begin transaction with error: %w", err)
	}

	defer func() {
		if rollErr := a.tran.Rollback(txnCtx); rollErr != nil {
			log.Printf("Error rolling back transaction: %v", rollErr)
		}
	}()

	accountTransactions, err := a.tr.GetAccountTransactions(accountId, getCtx)
	if err != nil {
		log.Printf("Unable to get transaction details for Account %s with error: %v", accountId, err)
		return nil, fmt.Errorf("unable to get transaction details with error: %w", err)
	}
	return accountTransactions, nil
}

func (a *AccountServiceImpl) Login(
	username string,
	password string,
	ctx context.Context,
) (*model.AccountDetails, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := a.tran.BeginTransaction(getCtx, transactional.IsolationLow, transactional.DurabilityLow)
	if err != nil {
		log.Printf("Error encountered when starting Login database transaction for "+
			"Username %s: ", username)
		return nil, fmt.Errorf("unable to begin transaction with error: %w", err)
	}

	defer func() {
		if rollErr := a.tran.Rollback(txnCtx); rollErr != nil {
			log.Printf("Error rolling back transaction: %v", rollErr)
		}
	}()

	accountDetails, err := a.ar.GetAccountDetailsFromUsername(username, getCtx)
	if err != nil {
		log.Printf("Unable to login with error: %v", err)
		if errors.Is(err, model.ErrNoMatchingUsername) {
			return nil, model.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("unable to login with error: %w", err)
	}
	exists := accountDetails.Password == password
	if !exists {
		log.Printf("Login failed for Username %s", username)
		return nil, model.ErrInvalidCredentials
	}
	return accountDetails, nil
}
