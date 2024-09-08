package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
	"webserver/migrations"
	"webserver/migrations/versions"
)

var accounts = []interface{}{
	mongodb.MongoAccountDetails{
		Id:       primitive.NewObjectID(),
		Username: "Olly",
		Password: "password",
		Person: mongodb.Person{
			FirstName: "Olly",
			LastName:  "OxenFree",
		},
		Accounts: []mongodb.Account{
			{
				Id:               primitive.NewObjectID(),
				AccountNumber:    "123-12345-0",
				AccountType:      "checking",
				AvailableBalance: 123.32,
			},
		},
		KnownAccounts: []mongodb.KnownAccount{},
		CreatedAt:     utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountDetails{
		Id:       primitive.NewObjectID(),
		Username: "Bob",
		Password: "bob'spassword",
		Person: mongodb.Person{
			FirstName: "Bob",
			LastName:  "Barker",
		},
		Accounts: []mongodb.Account{
			{
				Id:               primitive.NewObjectID(),
				AccountNumber:    "123-12345-1",
				AccountType:      "savings",
				AvailableBalance: 275.11,
			},
		},
		KnownAccounts: []mongodb.KnownAccount{
			{
				Id:            primitive.NewObjectID(),
				AccountNumber: "123-12345-0",
				AccountHolder: "Olly OxenFree",
				AccountType:   "checking",
			},
			{
				Id:            primitive.NewObjectID(),
				AccountNumber: "123-12345-2",
				AccountHolder: "Hilda Hill",
				AccountType:   "savings",
			},
		},
		CreatedAt: utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountDetails{
		Id:       primitive.NewObjectID(),
		Username: "Hilda",
		Password: "Hilda",
		Person: mongodb.Person{
			FirstName: "Hilda",
			LastName:  "Hill",
		},
		Accounts: []mongodb.Account{
			{
				Id:               primitive.NewObjectID(),
				AccountNumber:    "123-12345-2",
				AccountType:      "savings",
				AvailableBalance: 1004.55,
			},
			{
				Id:               primitive.NewObjectID(),
				AccountNumber:    "123-12345-3",
				AccountType:      "checking",
				AvailableBalance: 100.00,
			},
		},
		KnownAccounts: []mongodb.KnownAccount{
			{
				Id:            primitive.NewObjectID(),
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
		mongoCtx, cancel := context.WithTimeout(ctx, migrations.MigrationTimeout)
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
		mongoCtx, cancel := context.WithTimeout(ctx, migrations.MigrationTimeout)
		defer cancel()
		collectionName := "account"
		coll := db.Collection(collectionName)

		var ids = make([]primitive.ObjectID, len(accounts))
		for i, elem := range accounts {
			if mongoDetails, idOk := elem.(mongodb.MongoAccountDetails); idOk {
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
