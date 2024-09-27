package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"webserver/migrations/versions"
)

const collectionName = "migrations"
const MigrationTimeout = time.Minute * 1

type MigrationServiceImpl struct {
	client                *mongo.Client
	ctx                   context.Context
	migrationDatabaseName string
}

func NewMigrationService(client *mongo.Client, ctx context.Context, migrationDatabaseName string) *MigrationServiceImpl {
	return &MigrationServiceImpl{
		client:                client,
		ctx:                   ctx,
		migrationDatabaseName: migrationDatabaseName,
	}
}

// CheckIfApplied TODO: Set this function to private after finding out a way to test them
func (ms *MigrationServiceImpl) CheckIfApplied(version string) (bool, error) {
	db := ms.client.Database(ms.migrationDatabaseName)
	collection := db.Collection(collectionName)
	mongoCtx, cancel := context.WithTimeout(ms.ctx, MigrationTimeout)
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

// MarkAsApplied TODO: Set this function to private after finding out a way to test them
func (ms *MigrationServiceImpl) MarkAsApplied(version string) error {
	db := ms.client.Database(ms.migrationDatabaseName)
	collection := db.Collection(collectionName)
	mongoCtx, cancel := context.WithTimeout(ms.ctx, MigrationTimeout)
	defer cancel()
	_, err := collection.InsertOne(mongoCtx, bson.M{"version": version})
	if err != nil {
		return err
	}
	return nil
}

func (ms *MigrationServiceImpl) ApplyMigration(databaseName string, migration versions.Migration) (error, bool) {
	hasBeenApplied, err := ms.CheckIfApplied(migration.Version)
	if err != nil {
		log.Printf("Error when checking if migration %s has been applied: %v", migration.Version, err)
		return err, false
	}
	if hasBeenApplied {
		log.Printf("Migration %s has already been applied", migration.Version)
		return nil, false
	}
	err = migration.Up(ms.client, ms.ctx, databaseName)
	if err != nil {
		log.Printf("Error when applying migration %s: %v", migration.Version, err)
		return err, false
	}
	err = ms.MarkAsApplied(migration.Version)
	if err != nil {
		log.Printf("Error when marking migration %s as applied: %v", migration.Version, err)
		return err, true
	}
	return nil, true
}
