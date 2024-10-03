package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
	"webserver/migrations/versions"
)

const timeout = time.Minute * 1

var accountIds = []primitive.ObjectID{
	primitive.NewObjectID(),
	primitive.NewObjectID(),
	primitive.NewObjectID(),
	primitive.NewObjectID(),
}

var ollyAmount, _ = primitive.ParseDecimal128("123.23")
var bobAmount, _ = primitive.ParseDecimal128("275.11")
var hildaAmount1, _ = primitive.ParseDecimal128("1004.55")
var hildaAmount2, _ = primitive.ParseDecimal128("100.00")

var accounts = []interface{}{
	mongodb.MongoAccountOutput{
		Id:       primitive.NewObjectID(),
		Username: "Olly",
		Password: "password",
		Person: mongodb.Person{
			FirstName: "Olly",
			LastName:  "OxenFree",
		},
		BankAccounts: []mongodb.BankAccount{
			{
				Id:               accountIds[0],
				AccountNumber:    "123-12345-0",
				AccountType:      "checking",
				PendingBalance:   ollyAmount,
				AvailableBalance: ollyAmount,
			},
		},
		KnownBankAccounts: []mongodb.KnownBankAccount{},
		CreatedAt:         utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountOutput{
		Id:       primitive.NewObjectID(),
		Username: "Bob",
		Password: "bob'spassword",
		Person: mongodb.Person{
			FirstName: "Bob",
			LastName:  "Barker",
		},
		BankAccounts: []mongodb.BankAccount{
			{
				Id:               accountIds[1],
				AccountNumber:    "123-12345-1",
				AccountType:      "savings",
				PendingBalance:   bobAmount,
				AvailableBalance: bobAmount,
			},
		},
		KnownBankAccounts: []mongodb.KnownBankAccount{
			{
				Id:            accountIds[0],
				AccountNumber: "123-12345-0",
				AccountHolder: "Olly OxenFree",
				AccountType:   "checking",
			},
			{
				Id:            accountIds[2],
				AccountNumber: "123-12345-2",
				AccountHolder: "Hilda Hill",
				AccountType:   "savings",
			},
		},
		CreatedAt: utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountOutput{
		Id:       primitive.NewObjectID(),
		Username: "Hilda",
		Password: "Hilda",
		Person: mongodb.Person{
			FirstName: "Hilda",
			LastName:  "Hill",
		},
		BankAccounts: []mongodb.BankAccount{
			{
				Id:               accountIds[2],
				AccountNumber:    "123-12345-2",
				AccountType:      "savings",
				AvailableBalance: hildaAmount1,
				PendingBalance:   hildaAmount1,
			},
			{
				Id:               accountIds[3],
				AccountNumber:    "123-12345-3",
				AccountType:      "checking",
				AvailableBalance: hildaAmount2,
				PendingBalance:   hildaAmount2,
			},
		},
		KnownBankAccounts: []mongodb.KnownBankAccount{
			{
				Id:            accountIds[0],
				AccountNumber: "123-12345-0",
				AccountHolder: "Olly OxenFree",
				AccountType:   "checking",
			},
		},
		CreatedAt: utils.GetCurrentTimestamp(),
	},
}

var MigrationData1 = versions.Migration{
	Version: "1__Data",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		collectionName := "account"
		coll := db.Collection(collectionName)

		_, err := coll.InsertMany(mongoCtx, accounts)
		if err != nil {
			return err
		}

		log.Printf("account seed data successfully created")
		return nil
	},
	Down: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		collectionName := "account"
		coll := db.Collection(collectionName)

		var ids = make([]primitive.ObjectID, len(accounts))
		for i, elem := range accounts {
			if mongoDetails, idOk := elem.(mongodb.MongoAccountOutput); idOk {
				ids[i] = mongoDetails.Id
			}
		}
		deleteFilter := bson.M{"_id": bson.M{"$in": ids}}
		_, err := coll.DeleteMany(mongoCtx, deleteFilter)
		if err != nil {
			return err
		}

		return nil
	},
}
