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
		CreatedAt:        utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountDetails{
		Id:               primitive.NewObjectID(),
		AvailableBalance: 275.11,
		CreatedAt:        utils.GetCurrentTimestamp(),
	},
	mongodb.MongoAccountDetails{
		Id:               primitive.NewObjectID(),
		AvailableBalance: 1004.55,
		CreatedAt:        utils.GetCurrentTimestamp(),
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
