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

const AccountCollectionName = "account"

var MigrationSchema1 = versions.Migration{
	Version: "1__Schema",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, service.MigrationTimeout)
		defer cancel()
		collection := AccountCollectionName
		validator := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{
					"username", "password", "person", "_createdAt",
				},
				"properties": bson.M{
					"username": bson.M{
						"bsonType":    "string",
						"description": "Username for the Account [required]",
					},
					"password": bson.M{
						"bsonType":    "string",
						"description": "Password for the Account [required]",
					},
					"person": bson.M{
						"bsonType": "object",
						"required": []string{"firstName", "lastName"},
						"properties": bson.M{
							"firstName": bson.M{
								"bsonType":    "string",
								"description": "First Name of the Account Holder [required]",
							},
							"lastName": bson.M{
								"bsonType":    "string",
								"description": "Last Name of the Account Holder [required]",
							},
						},
					},
					// Embedded Bank BankAccounts (Main Bank BankAccounts for the User)
					"bankAccounts": bson.M{
						"bsonType": "array",
						"items": bson.M{
							"bsonType": "object",
							"required": []string{"accountNumber", "accountType", "availableBalance"},
							"properties": bson.M{
								"accountNumber": bson.M{
									"bsonType":    "string",
									"description": "Account Number for the Bank Account [required]",
								},
								"accountType": bson.M{
									"bsonType":    "string",
									"description": "Bank Account Type [required]",
								},
								"availableBalance": bson.M{
									"bsonType":    "decimal",
									"description": "Available Balance for the Bank Account [required]",
								},
								"_id": bson.M{
									"bsonType":    "objectId",
									"description": "Unique ID for this Bank Account [optional]",
								},
							},
						},
					},
					// Known Bank BankAccounts (External or Third-Party BankAccounts)
					"knownBankAccounts": bson.M{
						"bsonType": "array",
						"items": bson.M{
							"bsonType": "object",
							"required": []string{"accountNumber", "accountHolder", "accountType"},
							"properties": bson.M{
								"accountNumber": bson.M{
									"bsonType":    "string",
									"description": "Account Number of the Known Bank Account [required]",
								},
								"accountHolder": bson.M{
									"bsonType":    "string",
									"description": "Account Holder Name of the Known Bank Account [required]",
								},
								"accountType": bson.M{
									"bsonType":    "string",
									"description": "Type of the Known Bank Account [required]",
								},
								// Optional: Include ID if necessary for future references
								"_id": bson.M{
									"bsonType":    "objectId",
									"description": "Unique ID for this Known Bank Account [optional]",
								},
							},
						},
					},
					"_createdAt": bson.M{
						"bsonType":    "timestamp",
						"description": "Date of Account Creation [required]",
					},
				},
			},
		}

		opts := options.CreateCollection().SetValidator(validator).SetValidationLevel("strict")

		err := db.CreateCollection(mongoCtx, collection, opts)
		if err != nil {
			return err
		}

		log.Printf("Collection %s created with validator rules", collection)
		return nil
	},
	Down: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, service.MigrationTimeout)
		defer cancel()
		collection := "Account"
		err := db.Collection(collection).Drop(mongoCtx)
		if err != nil {
			return err
		}
		return nil
	},
}
