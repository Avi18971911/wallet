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
				"required": []string{"availableBalance", "_createdAt"},
				"properties": bson.M{
					"availableBalance": bson.M{
						"bsonType":    "double",
						"minimum":     0,
						"description": "must be a long and is required",
					},
					"_createdAt": bson.M{
						"bsonType":    "timestamp",
						"description": "must be a timestamp and is required",
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
