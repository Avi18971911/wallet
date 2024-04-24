//go:build test

package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"webserver/internal/pkg/infrastructure/mongodb/migrations/versions"
)

func ApplyMigrations(
	client *mongo.Client,
	ctx context.Context,
	mainDatabaseName string,
	migrationDatabaseName string,
	migrationsToRun []versions.Migration,
) {
	applyMigrations(
		client,
		ctx,
		mainDatabaseName,
		migrationDatabaseName,
		migrationsToRun,
	)
}

func CheckIfApplied(client *mongo.Client, ctx context.Context, databaseName string, version string) (bool, error) {
	return checkIfApplied(client, ctx, databaseName, version)
}

func MarkAsApplied(client *mongo.Client, ctx context.Context, databaseName string, version string) error {
	return markAsApplied(client, ctx, databaseName, version)
}
