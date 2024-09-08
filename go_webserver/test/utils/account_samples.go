package utils

import (
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
)

var TomAccountDetails = mongodb.MongoAccountInput{
	Username: "Tom",
	Password: "pass",
	Accounts: []mongodb.Account{
		{
			AccountNumber:    "123-45678-9",
			AccountType:      "savings",
			AvailableBalance: 1000,
		},
	},
	Person: mongodb.Person{
		FirstName: "Tom",
		LastName:  "Smith",
	},
	KnownAccounts: []mongodb.KnownAccount{},
	CreatedAt:     pkgutils.GetCurrentTimestamp(),
}

var SamAccountDetails = mongodb.MongoAccountInput{
	Username: "Sam",
	Password: "word",
	Accounts: []mongodb.Account{
		{
			AccountNumber:    "987-65432-1",
			AccountType:      "checking",
			AvailableBalance: 500,
		},
	},
	Person: mongodb.Person{
		FirstName: "Sam",
		LastName:  "Jones",
	},
	KnownAccounts: []mongodb.KnownAccount{},
	CreatedAt:     pkgutils.GetCurrentTimestamp(),
}
