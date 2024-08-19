package utils

import (
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
)

var TomAccountDetails = mongodb.MongoAccountInput{
	Username:        "Tom",
	Password:        "pass",
	AccountNumber:   "1234567890",
	AccountType:     "Savings",
	StartingBalance: 1000,
	Person: mongodb.Person{
		FirstName: "Tom",
		LastName:  "Smith",
	},
	KnownAccounts: []mongodb.KnownAccount{},
	CreatedAt:     pkgutils.GetCurrentTimestamp(),
}

var SamAccountDetails = mongodb.MongoAccountInput{
	Username:        "Sam",
	Password:        "word",
	AccountNumber:   "0987654321",
	AccountType:     "Savings",
	StartingBalance: 1000,
	Person: mongodb.Person{
		FirstName: "Sam",
		LastName:  "Jones",
	},
	KnownAccounts: []mongodb.KnownAccount{},
	CreatedAt:     pkgutils.GetCurrentTimestamp(),
}
