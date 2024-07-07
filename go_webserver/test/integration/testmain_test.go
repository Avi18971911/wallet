//go:build test

package integration

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"testing"
	"time"
	"webserver/test/utils"
)

var mongoClient *mongo.Client = nil
var cleanup func()

func TestMain(m *testing.M) {
	// Global setup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var err error
	mongoClient, cleanup = utils.CreateMongoRuntime(ctx)
	mongoURI := "mongodb://mongo:30001/?replicaSet=rs0"
	if err != nil {
		log.Fatalf("Failed to set up MongoDB runtime: %v", err)
	}

	err = utils.StartMigrationsContainer(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	// Run tests
	code := m.Run()

	// Global teardown
	cleanup()

	os.Exit(code)
}
