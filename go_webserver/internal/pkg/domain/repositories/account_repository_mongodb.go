package repositories

import (
	"context"
	"errors"
	"fmt"
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
		return nil, fmt.Errorf("error when converting account ID to object ID for accountId %s: %w", accountId, err)
	}
	filter := bson.M{"_id": objectId}
	err = ar.col.FindOne(ctx, filter).Decode(&accountDetails)
	if err != nil {
		return nil, fmt.Errorf("error when finding account by ID %s: %w", accountId, err)
	}
	res, err = fromMongoAccountDetails(&accountDetails)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully retrieved account details for account %s\n", accountId)
	return res, nil
}

func (ar *AccountRepositoryMongodb) AddBalance(
	accountId string,
	amount float64,
	ctx context.Context,
) error {
	objectId, err := utils.StringToObjectId(accountId)
	if err != nil {
		return fmt.Errorf("error when converting account ID to object ID for accountId %s: %w", accountId, err)
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{"$inc": bson.M{"availableBalance": amount}}
	result, err := ar.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(
			"error when updating account balance for accountId %s: %w", accountId, err,
		)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no matching account found for accountId %s", accountId)
	} else if result.ModifiedCount == 0 {
		return fmt.Errorf("update failed to the account balance for accountId %s", accountId)
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
		return fmt.Errorf("error when converting account ID to object ID for accountId %s: %w", accountId, err)
	}
	negativeAmount := amount * -1
	filter := bson.M{"_id": objectId}
	update := bson.M{"$inc": bson.M{"availableBalance": negativeAmount}}
	result, err := ar.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error when updating account balance for accountId %s: %w", accountId, err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no matching account found for accountId %s", accountId)
	} else if result.ModifiedCount == 0 {
		return fmt.Errorf("update failed to the account balance for accountId %s", accountId)
	} else {
		log.Printf("Successfully updated balance for account %s\n", accountId)
		return nil
	}
}

func (ar *AccountRepositoryMongodb) GetAccountDetailsFromUsername(
	username string,
	ctx context.Context,
) (*model.AccountDetails, error) {
	var accountDetails mongodb.MongoAccountDetails
	filter := bson.M{"username": username}
	err := ar.col.FindOne(ctx, filter).Decode(&accountDetails)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrNoMatchingUsername
		}
		return nil, fmt.Errorf("error when finding account by username: %s", err.Error())
	}
	return fromMongoAccountDetails(&accountDetails)
}

func fromMongoAccountDetails(details *mongodb.MongoAccountDetails) (*model.AccountDetails, error) {
	accountId, err := utils.ObjectIdToString(details.Id)
	if err != nil {
		return nil, fmt.Errorf(
			"error when converting object ID to string for username %s : %w", details.Username, err,
		)
	}
	return &model.AccountDetails{
		Id:               accountId,
		Username:         details.Username,
		Password:         details.Password,
		AvailableBalance: details.AvailableBalance,
		CreatedAt:        utils.TimestampToTime(details.CreatedAt),
	}, nil
}
