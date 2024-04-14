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

		collection := client.Database("test").Collection("testcol")
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
		assert.Equal(t, "value", result.Value)
	})
}
