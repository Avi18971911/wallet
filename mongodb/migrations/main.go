package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"migrations/versions"
	"migrations/versions/schema"
	"time"
)

// TODO: Change these to a file trawler
var migrationsToRun = []versions.Migration{
	schema.Migration1,
	schema.Migration2,
}

func main() {
	mainDatabaseName, migrationDatabaseName := "wallet", "migrations"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	// TODO: Set URI in config
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatalf("Error in connecting to database: %v", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error encountered when closing database connection: %v", err)
		}
	}(client, ctx)

	for _, elem := range migrationsToRun {
		hasBeenApplied, err := checkIfApplied(client, ctx, migrationDatabaseName, elem.Version)
		if err != nil {
			log.Fatalf("Error when checking if migration has been applied: %v", err)
		}
		if !hasBeenApplied {
			err = elem.Up(client, ctx, mainDatabaseName)
			if err != nil {
				log.Fatalf("Error when applying migration %s: %v", elem.Version, err)
			}
			err = markAsApplied(client, ctx, migrationDatabaseName, elem.Version)
			if err != nil {
				log.Printf("Error when marking migration %s as applied", elem.Version)
			}
		}
	}

	log.Println("Migration completed successfully")
}
