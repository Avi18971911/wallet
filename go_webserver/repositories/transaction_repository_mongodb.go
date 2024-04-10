package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/domain"
)

type TransactionRepositoryMongodb struct {
	col *mongo.Collection
}

func CreateNewTransactionRepositoryMongodb(col *mongo.Collection) *TransactionRepositoryMongodb {
	ar := TransactionRepositoryMongodb{col: col}
	return &ar
}

func (ar *TransactionRepositoryMongodb) AddTransaction(
	details domain.TransactionDetails,
	ctx context.Context,
) {
	_, err := ar.col.InsertOne(ctx, details)
	if err != nil {
		log.Fatalf("Failed to connect to Collection %s with error: %v", ar.col.Name(), err)
	}
	// TODO: Think about whether I want to return a bool or throw an exception potentially
}
