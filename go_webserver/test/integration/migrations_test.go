package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/migrations/service"
	"webserver/migrations/versions/schema"
	"webserver/test/utils"
)

const TestMigrationDatabaseName = "test_migrations"
const TestMigrationCollectionName = "test_migrations"

func TestMigrationService(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(TestMigrationDatabaseName)
	collection := db.Collection(TestMigrationCollectionName)
	ms := service.NewMigrationService(mongoClient, ctx, TestMigrationDatabaseName, TestMigrationCollectionName)

	t.Run("checkIfApplied should return true if a migration has already been applied", func(t *testing.T) {
		version := "1"
		_, err := collection.InsertOne(ctx, bson.M{"version": version})
		res, err := ms.CheckIfApplied(version)
		assert.Equal(t, true, res)
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})

	t.Run("checkIfApplied should return false if a migration hasn't been applied", func(t *testing.T) {
		version := "1"
		res, err := ms.CheckIfApplied(version)
		assert.Equal(t, false, res)
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})

	t.Run("checkIfApplied should return false if a migration hasn't been applied", func(t *testing.T) {
		version := "1"
		res, err := ms.CheckIfApplied(version)
		assert.Equal(t, false, res)
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})

	t.Run("markAsApplied should insert the migrated version into the database", func(t *testing.T) {
		version := "20"
		err := ms.MarkAsApplied(version)
		assert.Nil(t, err)
		err = collection.FindOne(ctx, bson.M{"version": version}).Err()
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})

	t.Run("checkIfApplied should return false after calling markAsApplied", func(t *testing.T) {
		version := "20"
		err := ms.MarkAsApplied(version)
		assert.Nil(t, err)
		res, err := ms.CheckIfApplied(version)
		assert.Equal(t, true, res)
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})
}

func TestV1SchemaMigration(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(TestMigrationDatabaseName)
	collection := db.Collection("account")
	migration := schema.MigrationSchema1
	utils.CleanupCollection(collection, ctx)

	t.Run("Should be able to add accounts with required fields", func(t *testing.T) {
		err := migration.Up(mongoClient, ctx, TestMigrationDatabaseName)
		assert.Nil(t, err)
		_, err = collection.InsertOne(ctx, sampleAccountDetails)
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})
}

func TestV2SchemaMigration(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	db := mongoClient.Database(TestMigrationDatabaseName)
	collection := db.Collection("transaction")
	migration := schema.MigrationSchema2
	utils.CleanupCollection(collection, ctx)

	t.Run("Should be able to add transactions with required fields", func(t *testing.T) {
		err := migration.Up(mongoClient, ctx, TestMigrationDatabaseName)
		assert.Nil(t, err)
		_, err = collection.InsertOne(ctx, sampleTransactionDetails)
		assert.Nil(t, err)
		utils.CleanupCollection(collection, ctx)
	})
}

var sampleAccountDetails = utils.TomAccountDetails

var sampleTransactionDetails = mongodb.MongoTransactionInput{
	FromBankAccountId: primitive.NewObjectID(),
	ToBankAccountId:   primitive.NewObjectID(),
	Amount:            primitive.NewDecimal128(1000, 0),
	Type:              "realized",
	CreatedAt:         pkgutils.GetCurrentTimestamp(),
}
