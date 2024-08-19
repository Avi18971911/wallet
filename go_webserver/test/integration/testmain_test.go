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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var err error
	mongoClient, cleanup = utils.CreateMongoRuntime(ctx)

	if err != nil {
		log.Fatalf("Failed to set up MongoDB runtime: %v", err)
	}
	code := m.Run()
	cleanup()
	os.Exit(code)
}
