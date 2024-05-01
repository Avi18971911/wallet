package migrations

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const collectionName = "migrations"
const MigrationTimeout = time.Minute * 1

func CheckIfApplied(client *mongo.Client, ctx context.Context, databaseName string, version string) (bool, error) {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	mongoCtx, cancel := context.WithTimeout(ctx, MigrationTimeout)
	defer cancel()
	filter := bson.M{"version": version}
	err := collection.FindOne(mongoCtx, filter).Err()
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
	mongoCtx, cancel := context.WithTimeout(ctx, MigrationTimeout)
	defer cancel()
	_, err := collection.InsertOne(mongoCtx, bson.M{"version": version})
	if err != nil {
		return err
	}
	return nil
}
