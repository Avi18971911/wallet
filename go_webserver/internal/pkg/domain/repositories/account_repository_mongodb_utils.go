package repositories

import (
	"fmt"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
)

func fromMongoAccountType(accountType string) int {
	switch accountType {
	case "savings":
		return model.Savings
	case "checking":
		return model.Checking
	case "investment":
		return model.Investment
	default:
		return -1
	}
}

func fromMongoKnownAccount(knownAccount []mongodb.KnownBankAccount) ([]model.KnownBankAccount, error) {
	var res = make([]model.KnownBankAccount, len(knownAccount))
	for i, ka := range knownAccount {
		stringId, err := utils.ObjectIdToString(ka.Id)
		if err != nil {
			return nil, fmt.Errorf(
				"error when converting object ID to string for known account %s: %v", ka.AccountNumber, err,
			)
		}
		res[i] = model.KnownBankAccount{
			Id:            stringId,
			AccountNumber: ka.AccountNumber,
			AccountHolder: ka.AccountHolder,
			AccountType:   fromMongoAccountType(ka.AccountType),
		}
	}
	return res, nil
}

func fromMongoAccounts(accounts []mongodb.BankAccount) ([]model.BankAccount, error) {
	var res = make([]model.BankAccount, len(accounts))
	for i, a := range accounts {
		stringId, err := utils.ObjectIdToString(a.Id)
		if err != nil {
			return nil, fmt.Errorf(
				"error when converting object ID to string for account number %s: %v", a.AccountNumber, err,
			)
		}
		availableBalanceDecimal, err := utils.FromPrimitiveDecimal128ToDecimal(a.AvailableBalance)
		if err != nil {
			return nil, fmt.Errorf(
				"error when converting available balance to decimal for account number %s: %v", a.AccountNumber, err,
			)
		}
		res[i] = model.BankAccount{
			Id:               stringId,
			AccountNumber:    a.AccountNumber,
			AccountType:      fromMongoAccountType(a.AccountType),
			AvailableBalance: availableBalanceDecimal,
		}
	}
	return res, nil
}

func fromMongoAccountDetails(details *mongodb.MongoAccountOutput) (*model.AccountDetails, error) {
	accountId, err := utils.ObjectIdToString(details.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"error when converting object ID to string for username %s : %w", details.Username, err,
		)
	}
	mongoAccounts, err := fromMongoAccounts(details.BankAccounts)
	if err != nil {
		return nil, fmt.Errorf("error when converting mongo accounts to model accounts: %v", err)
	}
	knownAccounts, err := fromMongoKnownAccount(details.KnownBankAccounts)
	if err != nil {
		return nil, fmt.Errorf("error when converting mongo known accounts to model known accounts: %v", err)
	}

	return &model.AccountDetails{
		Id:       accountId,
		Username: details.Username,
		Password: details.Password,
		Person: model.Person{
			FirstName: details.Person.FirstName,
			LastName:  details.Person.LastName,
		},
		BankAccounts:      mongoAccounts,
		KnownBankAccounts: knownAccounts,
		CreatedAt:         utils.TimestampToTime(details.CreatedAt),
	}, nil
}
