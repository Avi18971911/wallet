package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"regexp"
	"time"
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
	input *model.TransactionsForBankAccountInput,
	ctx context.Context,
) ([]model.BankAccountTransactionOutput, error) {
	getCtx, cancel := context.WithTimeout(ctx, addTimeout)
	defer cancel()

	txnCtx, err := a.tran.BeginTransaction(getCtx, transactional.IsolationHigh, transactional.DurabilityHigh)
	if err != nil {
		log.Printf("Error encountered when starting Get BankAccount Transactions database transaction for "+
			"BankAccount %s: %v", input.BankAccountId, err)
		return nil, fmt.Errorf("unable to begin transaction with error: %w", err)
	}

	defer func() {
		if rollErr := a.tran.Rollback(txnCtx); rollErr != nil {
			log.Printf("Error rolling back transaction: %v", rollErr)
		}
	}()

	accountTransactions, err := a.tr.GetTransactionsFromBankAccountId(input, getCtx)
	if err != nil {
		log.Printf(
			"Unable to get transaction details for BankAccount %s with error: %v",
			input.BankAccountId,
			err,
		)
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
	input *model.TransactionsForBankAccountInput,
	ctx context.Context,
) (model.AccountBalanceMonthsOutput, error) {
	transactions, err := a.GetBankAccountTransactions(input, ctx)
	defaultOutput := model.AccountBalanceMonthsOutput{}
	if err != nil {
		log.Printf("Unable to get Account History for BankAccount %s with error: %v", input.BankAccountId, err)
		return defaultOutput, fmt.Errorf("unable to get account history with error: %w", err)
	}
	availableBalance, pendingBalance, err := a.ar.GetAccountBalance(input.BankAccountId, ctx)
	if err != nil {
		log.Printf("Unable to get Account Balance for BankAccount %s with error: %v", input.BankAccountId, err)
		return defaultOutput, fmt.Errorf("unable to get account balance with error: %w", err)
	}
	months := getAccountBalanceMonths(
		input.BankAccountId,
		transactions,
		availableBalance,
		input.FromTime,
		input.ToTime,
		pendingBalance,
	)
	return months, nil
}

func getAccountBalanceMonths(
	bankAccountId string,
	transactions []model.BankAccountTransactionOutput,
	availableBalance decimal.Decimal,
	fromTime time.Time,
	toTime time.Time,
	pendingBalance decimal.Decimal,
) model.AccountBalanceMonthsOutput {
	numMonths := (toTime.Year()-fromTime.Year())*12 + int(toTime.Month()) - int(fromTime.Month()) + 1
	months := make([]model.AccountBalanceMonth, numMonths)
	latestMonth := model.AccountBalanceMonth{
		Month:            toTime.Month(),
		Year:             toTime.Year(),
		AvailableBalance: availableBalance,
		PendingBalance:   pendingBalance,
	}
	months[0] = latestMonth
	transactionsMonthMap := createTransactionsMonthMap(transactions)
	for i := 1; i < numMonths; i++ {
		previousMonthBalance := months[i-1]
		currentYear := previousMonthBalance.Year
		if previousMonthBalance.Month == time.January {
			currentYear--
		}
		currentMonth := time.December
		if previousMonthBalance.Month != time.January {
			currentMonth = previousMonthBalance.Month - 1
		}
		currentAccountBalanceMonth := model.AccountBalanceMonth{
			Month:            currentMonth,
			Year:             currentYear,
			AvailableBalance: previousMonthBalance.AvailableBalance,
			PendingBalance:   previousMonthBalance.PendingBalance,
		}
		key := int(currentMonth) + currentYear*12
		for _, transaction := range transactionsMonthMap[key] {
			currentAccountBalanceMonth.AvailableBalance, currentAccountBalanceMonth.PendingBalance = undoTransaction(
				transaction,
				currentAccountBalanceMonth.AvailableBalance,
				currentAccountBalanceMonth.PendingBalance,
			)
		}
		months[i] = currentAccountBalanceMonth
	}
	return model.AccountBalanceMonthsOutput{
		BankAccountId: bankAccountId,
		Months:        months,
	}
}

func createTransactionsMonthMap(
	transactions []model.BankAccountTransactionOutput,
) map[int][]model.BankAccountTransactionOutput {
	transactionsMonthMap := make(map[int][]model.BankAccountTransactionOutput)
	for _, transaction := range transactions {
		month := transaction.CreatedAt.Month()
		year := transaction.CreatedAt.Year()
		key := int(month) + year*12
		_, ok := transactionsMonthMap[key]
		if !ok {
			transactionsMonthMap[key] = make([]model.BankAccountTransactionOutput, 0)
		}
		transactionsMonthMap[key] = append(transactionsMonthMap[key], transaction)
	}
	return transactionsMonthMap
}

func undoTransaction(
	transaction model.BankAccountTransactionOutput,
	availableBalance decimal.Decimal,
	pendingBalance decimal.Decimal,
) (decimal.Decimal, decimal.Decimal) {
	adjustmentAmount := transaction.Amount
	if transaction.TransactionNature == model.Debit {
		adjustmentAmount = adjustmentAmount.Neg()
	}
	if transaction.TransactionType == model.Realized {
		availableBalance = availableBalance.Add(adjustmentAmount)
		pendingBalance = pendingBalance.Add(adjustmentAmount)
	} else {
		if transaction.Status == model.Active {
			pendingBalance = pendingBalance.Add(adjustmentAmount)
		}
		// Do nothing for revoked transactions or applied transactions
		// In the former case, the transaction's amount is no longer reflected in the pending balance
		// In the latter case, the transaction's amount is reflected in a realized transaction
	}
	return availableBalance, pendingBalance
}
