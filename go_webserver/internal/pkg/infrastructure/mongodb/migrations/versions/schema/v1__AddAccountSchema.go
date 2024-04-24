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

var Migration1 = versions.Migration{
	Version: "1",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
		defer cancel()
		collection := "Account"
		validation := bson.M{
			"validator": bson.M{
				"$jsonSchema": bson.M{
					"bsonType": "object",
					"required": []string{"availableBalance", "_createdAt"},
					"properties": bson.M{
						"availableBalance": bson.M{
							"bsonType":    "long",
							"description": "must be a long and is required",
						},
						"_createdAt": bson.M{
							"bsonType":    "timestamp",
							"description": "must be a timestamp and is required",
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
		collection := "Account"
		err := db.Collection(collection).Drop(mongoCtx)
		if err != nil {
			return err
		}
		return nil
	},
}
