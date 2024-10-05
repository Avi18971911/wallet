package utils

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
)

var tomBalanceDecimal128, _ = primitive.ParseDecimal128("231.95")
var samBalanceDecimal128, _ = primitive.ParseDecimal128("56.18")
var tomBalanceDecimal, _ = decimal.NewFromString("231.95")

var TomAccountDetails = mongodb.MongoAccountInput{
	Username: "Tom",
	Password: "pass",
	BankAccounts: []mongodb.BankAccount{
		{
			Id:               primitive.NewObjectID(),
			AccountNumber:    "123-45678-9",
			AccountType:      "savings",
			AvailableBalance: tomBalanceDecimal128,
			PendingBalance:   tomBalanceDecimal128,
		},
	},
	Person: mongodb.Person{
		FirstName: "Tom",
		LastName:  "Smith",
	},
	KnownBankAccounts: []mongodb.KnownBankAccount{
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
	BankAccounts: []mongodb.BankAccount{
		{
			Id:               primitive.NewObjectID(),
			AccountNumber:    "987-65432-1",
			AccountType:      "checking",
			AvailableBalance: samBalanceDecimal128,
			PendingBalance:   samBalanceDecimal128,
		},
	},
	Person: mongodb.Person{
		FirstName: "Sam",
		LastName:  "Jones",
	},
	KnownBankAccounts: []mongodb.KnownBankAccount{},
	CreatedAt:         pkgutils.GetCurrentTimestamp(),
}

var TomAccountDetailsModel = model.AccountDetailsOutput{
	Username: "Tom",
	BankAccounts: []model.BankAccount{
		{
			Id:               "UUID",
			AccountNumber:    "123-45678-9",
			AvailableBalance: tomBalanceDecimal,
			AccountType:      model.Checking,
		},
	},
	Person: model.Person{
		FirstName: "Tom",
		LastName:  "Smith",
	},
	KnownBankAccounts: []model.KnownBankAccount{
		{
			Id:            "UUID",
			AccountNumber: "987-65432-1",
			AccountHolder: "Sam Jones",
			AccountType:   model.Savings,
		},
	},
	CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
}
