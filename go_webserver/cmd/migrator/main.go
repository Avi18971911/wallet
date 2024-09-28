package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"strconv"
	"webserver/migrations/service"
	"webserver/migrations/versions"
	"webserver/migrations/versions/data"
	"webserver/migrations/versions/schema"
)

func main() {
	mainDatabaseName, migrationDatabaseName, migrationCollectionName := "wallet", "migrations", "migrations"
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

	schemaStartVer := parseEnvAsInt("SCHEMA_START_VER", 1)
	schemaEndVer := parseEnvAsInt("SCHEMA_END_VER", 2)

	dataStartVer := parseEnvAsInt("DATA_START_VER", 1)
	dataEndVer := parseEnvAsInt("DATA_END_VER", 1)

	ms := service.NewMigrationService(client, ctx, migrationDatabaseName, migrationCollectionName)

	log.Printf("Applying schema migrations to database %s", mainDatabaseName)
	applyMigrations(ms, mainDatabaseName, schema.SchemaMigrations, schemaStartVer, schemaEndVer)
	log.Printf("Applying data migrations to database %s", mainDatabaseName)
	applyMigrations(ms, mainDatabaseName, data.DataMigrations, dataStartVer, dataEndVer)
	log.Println("Migrations completed successfully")
}

func applyMigrations(
	ms *service.MigrationServiceImpl,
	mainDatabaseName string,
	migrationsToRun []versions.Migration,
	startVer int,
	endVer int,
) {
	migrationsToRun = migrationsToRun[startVer-1 : endVer]
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

func parseEnvAsInt(envVar string, defaultValue int) int {
	val := os.Getenv(envVar)
	if val == "" {
		return defaultValue
	}
	ret, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Error when parsing %s as int: %v", envVar, err)
	}
	return ret
}
