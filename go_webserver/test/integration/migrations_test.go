//go:build test

package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/migrations"
	"webserver/migrations/versions/schema"
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

func TestV1SchemaMigration(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(utils.TestDatabaseName)
	collection := db.Collection("account")
	migration := schema.MigrationSchema1

	// TODO: Add further tests once the schema has been finalized
	t.Run("Should be able to add accounts with balances and _createdAt", func(t *testing.T) {
		err := migration.Up(mongoClient, ctx, utils.TestDatabaseName)
		assert.Nil(t, err)
		mongoTimestamp := pkgutils.GetCurrentTimestamp()
		_, err = collection.InsertOne(ctx, bson.M{"_createdAt": mongoTimestamp, "availableBalance": 1000.00})
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})
}

func TestV2SchemaMigration(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(utils.TestDatabaseName)
	collection := db.Collection("transaction")
	migration := schema.MigrationSchema2

	// TODO: Add further tests once the schema has been finalized
	t.Run("Should be able to add transactions with required fields", func(t *testing.T) {
		err := migration.Up(mongoClient, ctx, utils.TestDatabaseName)
		assert.Nil(t, err)
		mongoTimestamp := pkgutils.GetCurrentTimestamp()
		accountId1, accountId2 := primitive.NewObjectID(), primitive.NewObjectID()
		_, err = collection.InsertOne(ctx, bson.M{
			"_createdAt":  mongoTimestamp,
			"amount":      1000.00,
			"fromAccount": accountId1,
			"toAccount":   accountId2,
		})
		assert.Nil(t, err)
		cleanupMigrations(collection, ctx)
	})
}

func cleanupMigrations(collection *mongo.Collection, ctx context.Context) {
	collection.Drop(ctx)
}
