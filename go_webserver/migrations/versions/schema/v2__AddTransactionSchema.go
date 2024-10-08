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

		realizedTransferValidation := bson.M{
			"bsonType": "object",
			"required": []string{"amount", "_createdAt", "fromBankAccountId", "toBankAccountId", "type", "_id"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType":    "objectId",
					"description": "the unique identifier for the transaction [required]",
				},
				"amount": bson.M{
					"bsonType":    "decimal",
					"description": "the amount transferred [required]",
				},
				"_createdAt": bson.M{
					"bsonType":    "timestamp",
					"description": "the time the transactions has been created [required]",
				},
				"fromBankAccountId": bson.M{
					"bsonType":    "objectId",
					"description": "the bank account ID from which the amount is coming [required]",
				},
				"toBankAccountId": bson.M{
					"bsonType":    "objectId",
					"description": "the bank account ID to which the amount is going [required]",
				},
				"type": bson.M{
					"bsonType":    "string",
					"description": "the type of transaction [required]",
					"enum":        []string{"realized"},
				},
			},
		}

		pendingTransferValidation := bson.M{
			"bsonType": "object",
			"required": []string{"amount", "_createdAt", "fromBankAccountId", "toBankAccountId", "type", "_id"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType":    "objectId",
					"description": "the unique identifier for the transaction [required]",
				},
				"amount": bson.M{
					"bsonType":    "decimal",
					"description": "the amount transferred [required]",
				},
				"_createdAt": bson.M{
					"bsonType":    "timestamp",
					"description": "the time the transactions has been created [required]",
				},
				"fromBankAccountId": bson.M{
					"bsonType":    "objectId",
					"description": "the bank account ID from which the amount is coming [required]",
				},
				"toBankAccountId": bson.M{
					"bsonType":    "objectId",
					"description": "the bank account ID to which the amount is going [required]",
				},
				"type": bson.M{
					"bsonType":    "string",
					"description": "the type of transaction [required]",
					"enum":        []string{"pending"},
				},
				"expirationDate": bson.M{
					"bsonType":    "timestamp",
					"description": "the exact moment the transaction will expire [required]",
				},
				"status": bson.M{
					"bsonType":    "string",
					"description": "the status of the transaction [required]",
					"enum":        []string{"active", "applied", "revoked"},
				},
			},
		}

		validation := bson.M{
			"$jsonSchema": bson.M{
				"oneOf": []bson.M{
					realizedTransferValidation,
					pendingTransferValidation,
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
