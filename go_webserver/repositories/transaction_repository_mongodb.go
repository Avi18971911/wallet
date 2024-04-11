package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
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
) error {
	_, err := ar.col.InsertOne(ctx, details)
	if err != nil {
		return err
	}
	return nil
}
