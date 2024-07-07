package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"webserver/migrations"
	"webserver/migrations/versions"
	"webserver/migrations/versions/data"
	"webserver/migrations/versions/schema"
)

// TODO: Change these to a file trawler
var migrationsToRun = []versions.Migration{
	schema.MigrationSchema1,
	schema.MigrationSchema2,
	data.MigrationData1,
}

func main() {
	mainDatabaseName, migrationDatabaseName := "wallet", "migrations"
	mongoURL := os.Getenv("MONGO_URL")
	log.Printf("Attempting to connect to Mongo URL %s", mongoURL)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Error in connecting to database: %v", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error encountered when closing database connection: %v", err)
		}
	}(client, ctx)

	applyMigrations(client, ctx, mainDatabaseName, migrationDatabaseName, migrationsToRun)
	log.Println("Migrations completed successfully")
}

func applyMigrations(
	client *mongo.Client,
	ctx context.Context,
	mainDatabaseName string,
	migrationDatabaseName string,
	migrationsToRun []versions.Migration,
) {
	for _, elem := range migrationsToRun {
		hasBeenApplied, err := migrations.CheckIfApplied(client, ctx, migrationDatabaseName, elem.Version)
		if err != nil {
			log.Fatalf("Error when checking if migration has been applied: %v", err)
		}
		if !hasBeenApplied {
			err = elem.Up(client, ctx, mainDatabaseName)
			if err != nil {
				log.Fatalf("Error when applying migration %s: %v", elem.Version, err)
			}
			err = migrations.MarkAsApplied(client, ctx, migrationDatabaseName, elem.Version)
			if err != nil {
				log.Printf("Error when marking migration %s as applied", elem.Version)
			}
		}
	}
}
