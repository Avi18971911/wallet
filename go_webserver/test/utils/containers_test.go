//go:build test

package utils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
	"time"
)

func TestMongoDBIntegration(t *testing.T) {
	t.Run("Should be able to start up Mongo container for testing purposes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		client, cleanup := CreateMongoRuntime(ctx)
		defer cleanup()

		collection := client.Database(TestDatabaseName).Collection("testcol")
		_, err := collection.InsertOne(ctx, bson.M{"name": "test", "value": "value"})
		if err != nil {
			t.Fatalf("Failed to insert document: %v", err)
		}

		var result struct {
			Name  string
			Value string
		}
		err = collection.FindOne(ctx, bson.M{"name": "test"}).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to find document: %v", err)
		}
		assert.Equal(t, "value", result.Value, "The value we inserted should be the value we found")
	})

	t.Run("Cleanup should successfully destroy every collection in the test database", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		client, cleanup := CreateMongoRuntime(ctx)
		defer cleanup()

		database := client.Database(TestDatabaseName)
		collection := database.Collection("testcol")
		_, err := collection.InsertOne(ctx, bson.M{"name": "test", "value": "value"})
		if err != nil {
			t.Fatalf("Failed to insert document: %v", err)
		}

		err = CleanupDatabase(client, ctx)
		assert.Nil(t, err)

		cursor, _ := collection.Find(ctx, bson.D{})
		assert.Equal(t, 0, cursor.RemainingBatchLength(), "Now the test collection should be empty")
		collections, listErr := database.ListCollectionNames(ctx, bson.D{})
		assert.Nil(t, listErr)
		assert.Equal(t, 0, len(collections), "There should be no collections after cleanup")
	})
}
