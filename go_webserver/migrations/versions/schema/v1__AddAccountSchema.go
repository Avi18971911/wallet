package schema

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"webserver/migrations"
	"webserver/migrations/versions"
)

var MigrationSchema1 = versions.Migration{
	Version: "1__Schema",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, migrations.MigrationTimeout)
		defer cancel()
		collection := "account"
		validator := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"availableBalance", "username", "password", "_createdAt"},
				"properties": bson.M{
					"availableBalance": bson.M{
						"bsonType":    "double",
						"minimum":     0,
						"description": "Available Balance for the Account [required]",
					},
					"password": bson.M{
						"bsonType":    "string",
						"description": "Password for the Account [required]",
					},
					"username": bson.M{
						"bsonType":    "string",
						"description": "Username for the Account [required]",
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
		mongoCtx, cancel := context.WithTimeout(ctx, migrations.MigrationTimeout)
		defer cancel()
		collection := "Account"
		err := db.Collection(collection).Drop(mongoCtx)
		if err != nil {
			return err
		}
		return nil
	},
}
