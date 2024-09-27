package service

import (
	"webserver/migrations/versions"
)

type MigrationService interface {
	// ApplyMigration /**
	// * ApplyMigration applies a migration to the database
	// * Returns an error if the migration fails, and a boolean indicating if the migration was applied
	ApplyMigration(databaseName string, migration versions.Migration) (error, bool)
}
