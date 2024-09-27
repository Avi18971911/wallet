package schema

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"webserver/migrations/service"
	"webserver/migrations/versions"
)

const TransactionCollectionName = "transaction"

var MigrationSchema2 = versions.Migration{
	Version: "2__Schema",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, service.MigrationTimeout)
		defer cancel()
		validation := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"amount", "_createdAt", "fromAccount", "toAccount"},
				"properties": bson.M{
					"amount": bson.M{
						"bsonType":    "decimal",
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
		}

		opts := options.CreateCollection().SetValidator(validation).SetValidationLevel("strict")

		err := db.CreateCollection(mongoCtx, TransactionCollectionName, opts)
		if err != nil {
			return err
		}

		log.Printf("Collection %s created with validation rules", TransactionCollectionName)
		return nil
	},
	Down: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, service.MigrationTimeout)
		defer cancel()
		err := db.Collection(TransactionCollectionName).Drop(mongoCtx)
		if err != nil {
			return err
		}
		return nil
	},
}
