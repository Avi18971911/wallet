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

func fromMongoKnownAccount(knownAccount []mongodb.KnownAccount) ([]model.KnownAccount, error) {
	var res = make([]model.KnownAccount, len(knownAccount))
	for i, ka := range knownAccount {
		stringId, err := utils.ObjectIdToString(ka.Id)
		if err != nil {
			return nil, fmt.Errorf(
				"error when converting object ID to string for known account %s: %v", ka.AccountNumber, err,
			)
		}
		res[i] = model.KnownAccount{
			Id:            stringId,
			AccountNumber: ka.AccountNumber,
			AccountHolder: ka.AccountHolder,
			AccountType:   fromMongoAccountType(ka.AccountType),
		}
	}
	return res, nil
}

func fromMongoAccounts(accounts []mongodb.Account) ([]model.Account, error) {
	var res = make([]model.Account, len(accounts))
	for i, a := range accounts {
		stringId, err := utils.ObjectIdToString(a.Id)
		if err != nil {
			return nil, fmt.Errorf(
				"error when converting object ID to string for account number %s: %v", a.AccountNumber, err,
			)
		}
		res[i] = model.Account{
			Id:               stringId,
			AccountNumber:    a.AccountNumber,
			AccountType:      fromMongoAccountType(a.AccountType),
			AvailableBalance: a.AvailableBalance,
		}
	}
	return res, nil
}

func fromMongoAccountDetails(details *mongodb.MongoAccountDetails) (*model.AccountDetails, error) {
	accountId, err := utils.ObjectIdToString(details.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"error when converting object ID to string for username %s : %w", details.Username, err,
		)
	}
	mongoAccounts, err := fromMongoAccounts(details.Accounts)
	if err != nil {
		return nil, fmt.Errorf("error when converting mongo accounts to model accounts: %v", err)
	}
	knownAccounts, err := fromMongoKnownAccount(details.KnownAccounts)
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
		Accounts:      mongoAccounts,
		KnownAccounts: knownAccounts,
		CreatedAt:     utils.TimestampToTime(details.CreatedAt),
	}, nil
}
