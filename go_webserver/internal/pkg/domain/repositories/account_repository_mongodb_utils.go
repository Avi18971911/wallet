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

func fromMongoKnownAccount(knownAccount []mongodb.KnownAccount) []model.KnownAccount {
	var res = make([]model.KnownAccount, len(knownAccount))
	for i, ka := range knownAccount {
		res[i] = model.KnownAccount{
			AccountNumber: ka.AccountNumber,
			AccountHolder: ka.AccountHolder,
			AccountType:   fromMongoAccountType(ka.AccountType),
		}
	}
	return res
}

func fromMongoAccounts(accounts []mongodb.Account) []model.Account {
	var res = make([]model.Account, len(accounts))
	for i, a := range accounts {
		res[i] = model.Account{
			AccountNumber:    a.AccountNumber,
			AccountType:      fromMongoAccountType(a.AccountType),
			AvailableBalance: a.AvailableBalance,
		}
	}
	return res
}

func fromMongoAccountDetails(details *mongodb.MongoAccountDetails) (*model.AccountDetails, error) {
	accountId, err := utils.ObjectIdToString(details.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"error when converting object ID to string for username %s : %w", details.Username, err,
		)
	}
	return &model.AccountDetails{
		Id:       accountId,
		Username: details.Username,
		Password: details.Password,
		Person: model.Person{
			FirstName: details.Person.FirstName,
			LastName:  details.Person.LastName,
		},
		Accounts:      fromMongoAccounts(details.Accounts),
		KnownAccounts: fromMongoKnownAccount(details.KnownAccounts),
		CreatedAt:     utils.TimestampToTime(details.CreatedAt),
	}, nil
}
