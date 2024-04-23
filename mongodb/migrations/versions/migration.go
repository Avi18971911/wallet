package versions

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	Version string
	Up      func(db *mongo.Client, ctx context.Context, databaseName string) error
	Down    func(db *mongo.Client, ctx context.Context, databaseName string) error
}
