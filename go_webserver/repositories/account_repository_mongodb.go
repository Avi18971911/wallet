package repositories

import (
	"context"
	"errors"
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

func (ar *AccountRepositoryMongodb) GetAccountDetails(accountId string, ctx context.Context) (*domain.AccountDetails, error) {
	var accountDetails domain.AccountDetails
	filter := bson.M{"accountId": accountId}
	err := ar.col.FindOne(ctx, filter).Decode(&accountDetails)
	if err != nil {
		return nil, err
	}
	return &accountDetails, nil
}

func (ar *AccountRepositoryMongodb) GetAccountTransactions(
	accountId string, ctx context.Context,
) ([]*domain.AccountTransaction, error) {
	return []*domain.AccountTransaction{}, nil
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
		log.Printf("Successfully updated balance for account %s\n", accountId)
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
