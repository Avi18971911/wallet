package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
)

var TomAccountDetails = mongodb.MongoAccountInput{
	Username: "Tom",
	Password: "pass",
	Accounts: []mongodb.Account{
		{
			Id:               primitive.NewObjectID(),
			AccountNumber:    "123-45678-9",
			AccountType:      "savings",
			AvailableBalance: 1000,
		},
	},
	Person: mongodb.Person{
		FirstName: "Tom",
		LastName:  "Smith",
	},
	KnownAccounts: []mongodb.KnownAccount{
		{
			Id:            primitive.NewObjectID(),
			AccountNumber: "987-65432-1",
			AccountHolder: "Sam Jones",
			AccountType:   "checking",
		},
	},
	CreatedAt: pkgutils.GetCurrentTimestamp(),
}

var SamAccountDetails = mongodb.MongoAccountInput{
	Username: "Sam",
	Password: "word",
	Accounts: []mongodb.Account{
		{
			Id:               primitive.NewObjectID(),
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

var TomAccountDetailsModel = model.AccountDetails{
	Username: "Tom",
	Accounts: []model.Account{
		{
			Id:               "UUID",
			AccountNumber:    "123-45678-9",
			AvailableBalance: 1000,
			AccountType:      1,
		},
	},
	Person: model.Person{
		FirstName: "Tom",
		LastName:  "Smith",
	},
	KnownAccounts: []model.KnownAccount{
		{
			Id:            "UUID",
			AccountNumber: "987-65432-1",
			AccountHolder: "Sam Jones",
			AccountType:   0,
		},
	},
	CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
}
