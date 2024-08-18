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
		Id:               primitive.NewObjectID(),
		AvailableBalance: 123.32,
		Username:         "Olly",
		Password:         "password",
		AccountNumber:    "123-12345-0",
		Person: mongodb.Person{
			FirstName: "Olly",
			LastName:  "OxenFree",
		},
		AccountType: "Checking",
		CreatedAt:   utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountDetails{
		Id:               primitive.NewObjectID(),
		AvailableBalance: 275.11,
		Username:         "Bob",
		Password:         "bob'spassword",
		AccountNumber:    "123-12345-1",
		AccountType:      "Savings",
		Person: mongodb.Person{
			FirstName: "Bob",
			LastName:  "Barker",
		},
		KnownAccounts: []mongodb.KnownAccount{
			{
				AccountNumber: "123-12345-0",
				AccountHolder: "Olly OxenFree",
				AccountType:   "Checking",
			},
			{
				AccountNumber: "123-12345-2",
				AccountHolder: "Hilda Hill",
				AccountType:   "Savings",
			},
		},
		CreatedAt: utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountDetails{
		Id:               primitive.NewObjectID(),
		AvailableBalance: 1004.55,
		Username:         "Hilda",
		Password:         "Hilda",
		AccountNumber:    "123-12345-2",
		AccountType:      "Savings",
		Person: mongodb.Person{
			FirstName: "Hilda",
			LastName:  "Hill",
		},
		KnownAccounts: []mongodb.KnownAccount{
			{
				AccountNumber: "123-12345-0",
				AccountHolder: "Olly OxenFree",
				AccountType:   "Checking",
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

		log.Printf("inserting account seed data %s", accounts)
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
