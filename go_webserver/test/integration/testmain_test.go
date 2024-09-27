package integration

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"testing"
	"time"
	"webserver/migrations/service"
	"webserver/migrations/versions/schema"
	"webserver/test/utils"
)

var mongoClient *mongo.Client = nil
var cleanup func()

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var err error
	mongoClient, cleanup = utils.CreateMongoRuntime(ctx)

	if err != nil {
		log.Fatalf("Failed to set up MongoDB runtime: %v", err)
	}

	mainDatabaseName, migrationDatabaseName := utils.TestDatabaseName, "migrations"
	if mongoClient == nil {
		log.Fatalf("mongoClient is uninitialized or otherwise nil")
	}
	ms := service.NewMigrationService(mongoClient, ctx, migrationDatabaseName)
	migrations := schema.SchemaMigrations
	for _, elem := range migrations {
		_, hasBeenApplied := ms.ApplyMigration(mainDatabaseName, elem)
		if !hasBeenApplied {
			log.Fatalf("Unable to apply migration %s", elem.Version)
		}
	}

	code := m.Run()
	cleanup()
	os.Exit(code)
}
