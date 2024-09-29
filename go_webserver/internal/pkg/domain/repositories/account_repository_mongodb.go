package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (ar *AccountRepositoryMongodb) GetAccountDetailsFromBankAccountId(
	bankAccountId string,
	ctx context.Context,
) (*model.AccountDetails, error) {
	var accountDetails mongodb.MongoAccountOutput
	var res *model.AccountDetails
	objectId, err := utils.StringToObjectId(bankAccountId)
	if err != nil {
		return nil, fmt.Errorf("error when converting account ID to object ID for "+
			"bankAccountId %s: %w", bankAccountId, err)
	}
	filter := bson.M{"bankAccounts._id": objectId}
	err = ar.col.FindOne(ctx, filter).Decode(&accountDetails)
	if err != nil {
		return nil, fmt.Errorf("error when finding account by ID %s: %w", bankAccountId, err)
	}
	res, err = fromMongoAccountDetails(&accountDetails)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully retrieved account details for account %s\n", bankAccountId)
	return res, nil
}

func (ar *AccountRepositoryMongodb) AddBalance(
	bankAccountId string,
	amount decimal.Decimal,
	ctx context.Context,
) error {
	objectId, err := utils.StringToObjectId(bankAccountId)
	if err != nil {
		return fmt.Errorf("error when converting account ID to object ID for "+
			"bankAccountId %s: %w", bankAccountId, err)
	}
	decimal128Amount, err := utils.FromDecimalToPrimitiveDecimal128(amount)
	if err != nil {
		return fmt.Errorf("error when converting amount to Decimal128 for "+
			"bankAccountId %s: %w", bankAccountId, err)
	}
	filter := bson.M{"bankAccounts._id": objectId}
	update := bson.M{"$inc": bson.M{"bankAccounts.$.availableBalance": decimal128Amount}}
	result, err := ar.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf(
			"error when updating account balance for bankAccountId %s: %w", bankAccountId, err,
		)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("no matching account found for bankAccountId %s", bankAccountId)
	} else if result.ModifiedCount == 0 {
		return fmt.Errorf("update failed to the account balance for bankAccountId %s", bankAccountId)
	} else {
		log.Printf("Successfully updated balance for account %s\n", bankAccountId)
		return nil
	}
}

func (ar *AccountRepositoryMongodb) DeductBalance(
	bankAccountId string,
	amount decimal.Decimal,
	ctx context.Context,
) (decimal.Decimal, error) {
	objectId, err := utils.StringToObjectId(bankAccountId)
	defaultDecimal := decimal.NewFromInt(0)
	if err != nil {
		return defaultDecimal,
			fmt.Errorf("error when converting account ID to object ID for bankAccountId "+
				"%s: %w", bankAccountId, err)
	}
	negativeAmount := amount.Neg()
	decimal128Amount, err := utils.FromDecimalToPrimitiveDecimal128(negativeAmount)
	if err != nil {
		return defaultDecimal,
			fmt.Errorf("error when converting amount to Decimal128 for bankAccountId %s: %w", bankAccountId, err)
	}
	filter := bson.M{"bankAccounts._id": objectId}
	update := bson.M{"$inc": bson.M{"bankAccounts.$.availableBalance": decimal128Amount}}
	result, err := ar.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return defaultDecimal,
			fmt.Errorf("error when updating account balance for bankAccountId %s: %w", bankAccountId, err)
	}
	if result.MatchedCount == 0 {
		return defaultDecimal,
			fmt.Errorf("no matching account found for bankAccountId %s", bankAccountId)
	} else if result.ModifiedCount == 0 {
		return defaultDecimal,
			fmt.Errorf("update failed to the account balance for bankAccountId %s", bankAccountId)
	} else {
		updatedBalance, err := ar.GetAccountBalance(bankAccountId, ctx)
		if err != nil {
			return defaultDecimal,
				fmt.Errorf("error when getting updated balance for account %s: %w", bankAccountId, err)
		}
		return updatedBalance, nil
	}
}

func (ar *AccountRepositoryMongodb) GetAccountBalance(
	bankAccountId string,
	ctx context.Context,
) (decimal.Decimal, error) {
	objectId, err := utils.StringToObjectId(bankAccountId)
	defaultDecimal := decimal.NewFromInt(0)
	if err != nil {
		return defaultDecimal,
			fmt.Errorf("error when converting account ID to object ID for bankAccountId "+
				"%s: %w", bankAccountId, err)
	}
	filter := bson.M{"bankAccounts._id": objectId}
	projection := bson.M{"bankAccounts.$": 1}
	var account mongodb.MongoAccountInput
	err = ar.col.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&account)
	if err != nil {
		return defaultDecimal,
			fmt.Errorf("error when finding account by ID %s: %w", bankAccountId, err)
	}
	balanceDecimal, err := utils.FromPrimitiveDecimal128ToDecimal(account.BankAccounts[0].AvailableBalance)
	if err != nil {
		return defaultDecimal,
			fmt.Errorf("error when converting available balance to decimal for bankAccountId "+
				"%s: %w", bankAccountId, err)
	}
	return balanceDecimal, nil
}

func (ar *AccountRepositoryMongodb) GetAccountDetailsFromUsername(
	username string,
	ctx context.Context,
) (*model.AccountDetails, error) {
	var accountDetails mongodb.MongoAccountOutput
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
