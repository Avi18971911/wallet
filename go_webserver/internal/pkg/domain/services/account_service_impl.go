package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
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

func validateAccountDetails(accountDetails *model.AccountDetailsOutput) error {
	err := validateAccountNumbers(accountDetails.BankAccounts)
	if err != nil {
		return fmt.Errorf("unable to validate account numbers with error: %w", err)
	}
	return nil
}

func validateAccountNumbers(accounts []model.BankAccount) error {
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

func (a *AccountServiceImpl) GetAccountDetailsFromBankAccountId(
	bankAccountId string,
	ctx context.Context,
) (*model.AccountDetailsOutput, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()
	accountDetails, err := a.ar.GetAccountDetailsFromBankAccountId(bankAccountId, getCtx)
	if err != nil {
		log.Printf("Unable to get account details for BankAccount %s with error: %v", bankAccountId, err)
		return nil, fmt.Errorf("unable to get account details with error: %w", err)
	}
	err = validateAccountDetails(accountDetails)
	if err != nil {
		log.Printf("Unable to successfully validate account details for BankAccount %s with error: "+
			"%v", bankAccountId, err)
		return nil, fmt.Errorf("unable to validate account details with error: %w", err)
	}
	return accountDetails, nil
}

func (a *AccountServiceImpl) GetBankAccountTransactions(
	bankAccountId string, ctx context.Context,
) ([]model.BankAccountTransactionOutput, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := a.tran.BeginTransaction(getCtx, transactional.IsolationHigh, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Get BankAccount Transactions database transaction for "+
			"BankAccount %s: %v", bankAccountId, err)
		return nil, fmt.Errorf("unable to begin transaction with error: %w", err)
	}

	defer func() {
		if rollErr := a.tran.Rollback(txnCtx); rollErr != nil {
			log.Printf("Error rolling back transaction: %v", rollErr)
		}
	}()

	accountTransactions, err := a.tr.GetTransactionsFromBankAccountId(bankAccountId, getCtx)
	if err != nil {
		log.Printf("Unable to get transaction details for BankAccount %s with error: %v", bankAccountId, err)
		return nil, fmt.Errorf("unable to get transaction details with error: %w", err)
	}
	return accountTransactions, nil
}

func (a *AccountServiceImpl) Login(
	username string,
	password string,
	ctx context.Context,
) (*model.AccountDetailsOutput, error) {
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

func (a *AccountServiceImpl) GetAccountHistoryInMonths(
	bankAccountId string,
	startMonth
	ctx context.Context,
) ([]model.AccountBalanceMonths, error) {
	transactions, err := a.GetBankAccountTransactions(bankAccountId, ctx)
	if err != nil {
		log.Printf("Unable to get Account History for BankAccount %s with error: %v", bankAccountId, err)
		return nil, fmt.Errorf("unable to get account history with error: %w", err)
	}
	availableBalance, pendingBalance, err := a.ar.GetAccountBalance(bankAccountId, ctx)
	if err != nil {
		log.Printf("Unable to get Account Balance for BankAccount %s with error: %v", bankAccountId, err)
		return nil, fmt.Errorf("unable to get account balance with error: %w", err)
	}
	months := getAccountBalanceMonths(transactions, availableBalance, pendingBalance)
	return months, nil
}

func getAccountBalanceMonths(
	transactions []model.BankAccountTransactionOutput,
	availableBalance decimal.Decimal,
	pendingBalance decimal.Decimal,
) []model.AccountBalanceMonthsOutput {
	months := make([]model.AccountBalanceMonthsOutput, 0)
	month := model.AccountBalanceMonthsOutput{}
	for _, transaction := range transactions {
		if month.Month == transaction.Date.Month() {
			month.Transactions
		}
	}
}
