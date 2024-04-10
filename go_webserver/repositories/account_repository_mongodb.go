package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/domain"
)

type DB interface {
	GetOne(ctx context.Context, accountId string)
}
type AccountRepositoryMongodb struct {
	col *mongo.Collection
}

func CreateNewAccountRepositoryMongodb(ctx context.Context, col *mongo.Collection) *AccountRepositoryMongodb {
	ar := AccountRepositoryMongodb{col: col}
	return &ar
}

func (ar *AccountRepositoryMongodb) GetAccountDetails(accountId string, ctx context.Context) *domain.AccountDetails {
	var accountDetails domain.AccountDetails
	record, err := ar.col.Find(ctx, bson.D{{"accountId", accountId}})
	if err != nil {
		log.Fatalf("Failed to connect to Collection %s with error: %v", ar.col.Name(), err)
	}
	err = record.Decode(&accountDetails)
	if err != nil {
		log.Fatalf("Failed to decode the Document %s with error: %v", record.ID(), err)
	}
	return &accountDetails
}

func (ar *AccountRepositoryMongodb) GetAccountTransactions(
	accountId string, ctx context.Context,
) []*domain.AccountTransaction {
	return []*domain.AccountTransaction{}
}
