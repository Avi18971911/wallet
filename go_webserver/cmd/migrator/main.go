package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"webserver/migrations/service"
	"webserver/migrations/versions"
	"webserver/migrations/versions/data"
	"webserver/migrations/versions/schema"
)

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

	ms := service.NewMigrationService(client, ctx, migrationDatabaseName)
	applyMigrations(ms, mainDatabaseName, migrationsToRun)
	log.Println("Migrations completed successfully")
}

func applyMigrations(
	ms *service.MigrationServiceImpl,
	mainDatabaseName string,
	migrationsToRun []versions.Migration,
) {
	for _, elem := range migrationsToRun {
		err, hasBeenApplied := ms.ApplyMigration(mainDatabaseName, elem)
		if err != nil {
			log.Fatalf("Error when applying migration %s: %v", elem.Version, err)
		}
		if hasBeenApplied {
			log.Printf("Migration %s has been applied", elem.Version)
		} else {
			log.Printf("Migration %s has not been applied", elem.Version)
		}
	}
}
