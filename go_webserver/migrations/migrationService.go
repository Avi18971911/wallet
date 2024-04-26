package migrations

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "migrations"

func CheckIfApplied(client *mongo.Client, ctx context.Context, databaseName string, version string) (bool, error) {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	filter := bson.M{"version": version}
	err := collection.FindOne(ctx, filter).Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
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
