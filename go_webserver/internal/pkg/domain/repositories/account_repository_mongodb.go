package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
)

type AccountRepositoryMongodb struct {
	col *mongo.Collection
}

func CreateNewAccountRepositoryMongodb(col *mongo.Collection) *AccountRepositoryMongodb {
	ar := AccountRepositoryMongodb{col: col}
	return &ar
}

func (ar *AccountRepositoryMongodb) GetAccountDetails(
	accountId string,
	ctx context.Context,
) (*model.AccountDetails, error) {
	var accountDetails mongodb.MongoAccountDetails
	var res *model.AccountDetails
	objectId, err := utils.StringToObjectId(accountId)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectId}
	err = ar.col.FindOne(ctx, filter).Decode(&accountDetails)
	if err != nil {
		return nil, err
	}
	res, err = fromMongoAccountDetails(&accountDetails)
	return res, nil
}

func (ar *AccountRepositoryMongodb) AddBalance(
	accountId string,
	amount float64,
	ctx context.Context,
) error {
	objectId, err := utils.StringToObjectId(accountId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectId}
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
	objectId, err := utils.StringToObjectId(accountId)
	if err != nil {
		return err
	}
	negativeAmount := amount * -1
	filter := bson.M{"_id": objectId}
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

func fromMongoAccountDetails(details *mongodb.MongoAccountDetails) (*model.AccountDetails, error) {
	accountId, err := utils.ObjectIdToString(details.Id)
	if err != nil {
		return nil, err
	}
	return &model.AccountDetails{
		Id:               accountId,
		Username:         details.Username,
		AvailableBalance: details.AvailableBalance,
		CreatedAt:        details.CreatedAt,
	}, nil
}
