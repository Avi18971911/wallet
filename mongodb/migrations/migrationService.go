package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "migrations"

func CheckIfApplied(client *mongo.Client, ctx context.Context, databaseName string, version string) bool {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	filter := bson.M{"version": version}
	err := collection.FindOne(ctx, filter).Err()
	return err != nil
}

func MarkAsApplied(client *mongo.Client, ctx context.Context, databaseName string, version string) error {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	_, err := collection.InsertOne(ctx, bson.M{"version": version})
	if err != nil {
		return err
	}
	return nil
}
