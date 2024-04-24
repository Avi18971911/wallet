//go:build test

package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	"webserver/internal/pkg/infrastructure/mongodb/migrations"
	"webserver/internal/pkg/infrastructure/mongodb/migrations/versions/schema"
	"webserver/test/utils"
)

func TestMigrationService(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(utils.TestDatabaseName)
	collection := db.Collection("migrations")

	t.Run("checkIfApplied should return true if a migration has already been applied", func(t *testing.T) {
		version := "1"
		_, err := collection.InsertOne(ctx, bson.M{"version": version})
		res, err := migrations.CheckIfApplied(mongoClient, ctx, utils.TestDatabaseName, version)
		assert.Equal(t, true, res)
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})

	t.Run("checkIfApplied should return false if a migration hasn't been applied", func(t *testing.T) {
		version := "1"
		res, err := migrations.CheckIfApplied(mongoClient, ctx, utils.TestDatabaseName, version)
		assert.Equal(t, false, res)
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})

	t.Run("checkIfApplied should return false if a migration hasn't been applied", func(t *testing.T) {
		version := "1"
		res, err := migrations.CheckIfApplied(mongoClient, ctx, utils.TestDatabaseName, version)
		assert.Equal(t, false, res)
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})

	t.Run("markAsApplied should insert the migrated version into the database", func(t *testing.T) {
		version := "20"
		err := migrations.MarkAsApplied(mongoClient, ctx, utils.TestDatabaseName, version)
		assert.Nil(t, err)
		err = collection.FindOne(ctx, bson.M{"version": version}).Err()
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})

	t.Run("checkIfApplied should return false after calling markAsApplied", func(t *testing.T) {
		version := "20"
		err := migrations.MarkAsApplied(mongoClient, ctx, utils.TestDatabaseName, version)
		assert.Nil(t, err)
		res, err := migrations.CheckIfApplied(mongoClient, ctx, utils.TestDatabaseName, version)
		assert.Equal(t, true, res)
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})
}

func TestV1Migration(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(utils.TestDatabaseName)
	collection := db.Collection("account")
	migration := schema.Migration1

	t.Run("")

}

func cleanupMigrations(collection *mongo.Collection, ctx context.Context) {
	collection.Drop(ctx)
}
