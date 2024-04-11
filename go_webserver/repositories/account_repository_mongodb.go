package repositories

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/domain"
)

type AccountRepositoryMongodb struct {
	col *mongo.Collection
}

func CreateNewAccountRepositoryMongodb(col *mongo.Collection) *AccountRepositoryMongodb {
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

func (ar *AccountRepositoryMongodb) AddBalance(
	accountId string,
	amount float64,
	ctx context.Context,
) error {
	filter := bson.M{"accountId": accountId}
	update := bson.M{"$inc": bson.M{"availableBalance": amount}}
	result, err := ar.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no matching account found")
	} else if result.ModifiedCount == 0 {
		return errors.New("no update made to the account balance")
	} else {
		fmt.Printf("Successfully updated balance for account %s\n", accountId)
		return nil
	}
}

func (ar *AccountRepositoryMongodb) DeductBalance(
	accountId string,
	amount float64,
	ctx context.Context,
) error {
	negativeAmount := amount * -1
	filter := bson.M{"accountId": accountId}
	update := bson.M{"$inc": bson.M{"availableBalance": negativeAmount}}
	result, err := ar.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no matching account found")
	} else if result.ModifiedCount == 0 {
		return errors.New("no update made to the account balance")
	} else {
		log.Printf("Successfully updated balance for account %s\n", accountId)
		return nil
	}
}
