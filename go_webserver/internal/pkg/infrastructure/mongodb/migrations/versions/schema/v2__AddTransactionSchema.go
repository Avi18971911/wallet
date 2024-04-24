package schema

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"webserver/internal/pkg/infrastructure/mongodb/migrations/versions"
)

const collection = "Transaction"

var Migration2 = versions.Migration{
	Version: "2",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
		defer cancel()
		validation := bson.M{
			"validator": bson.M{
				"$jsonSchema": bson.M{
					"bsonType": "object",
					"required": []string{"amount", "_createdAt", "fromAccount", "toAccount"},
					"properties": bson.M{
						"amount": bson.M{
							"bsonType":    "long",
							"description": "the amount transferred",
						},
						"_createdAt": bson.M{
							"bsonType":    "timestamp",
							"description": "the time the transactions has been created",
						},
						"fromAccount": bson.M{
							"bsonType":    "objectId",
							"description": "the account from which the amount is coming",
						},
						"toAccount": bson.M{
							"bsonType":    "objectId",
							"description": "the account to which the amount is going",
						},
					},
				},
			},
			"validationLevel": "strict",
		}

		opts := options.CreateCollection().SetValidator(validation)

		err := db.CreateCollection(mongoCtx, collection, opts)
		if err != nil {
			return err
		}

		log.Printf("Collection %s created with validation rules", collection)
		return nil
	},
	Down: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
		defer cancel()
		err := db.Collection(collection).Drop(mongoCtx)
		if err != nil {
			return err
		}
		return nil
	},
}
