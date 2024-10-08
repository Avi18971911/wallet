package integration

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/migrations/versions/schema"
	"webserver/test/utils"
)

func TestGetAccountTransactions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	tomAccountName, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.BankAccounts[0].Id)
	samAccountName, _ := pkgutils.ObjectIdToString(utils.SamAccountDetails.BankAccounts[0].Id)
	tomObjectId, _ := pkgutils.StringToObjectId(tomAccountName)
	samObjectId, _ := pkgutils.StringToObjectId(samAccountName)

	tranAmounts := []decimal.Decimal{
		decimal.NewFromFloat(100.0),
		decimal.NewFromFloat(200.0),
		decimal.NewFromFloat(300.0),
	}

	tranAmountsDecimal128 := make([]primitive.Decimal128, 3)
	for i, _ := range tranAmounts {
		tranAmountsDecimal128[i], _ = pkgutils.FromDecimalToPrimitiveDecimal128(tranAmounts[i])
	}

	transactionsInput := makeTransactionsInput(tomObjectId, samObjectId, tranAmountsDecimal128)

	t.Run(
		"Allows the insertion of transactions and the retrieval of all transactions from an account",
		func(t *testing.T) {
			setupGetAccountTransactionsTestCase(accCollection, tranCollection, ctx, t)
			_, err := tranCollection.InsertMany(ctx, transactionsInput)
			if err != nil {
				t.Errorf("Error inserting transactions: %v", err)
			}

			accountService := setupAccountService(mongoClient, tranCollection, accCollection)

			input := model.TransactionsForBankAccountInput{
				BankAccountId: tomAccountName,
				FromTime:      time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
				ToTime:        time.Date(2050, time.December, 31, 23, 59, 59, 0, time.UTC),
			}
			res, _ := accountService.GetBankAccountTransactions(&input, ctx)
			expectedCreatedAts := []time.Time{firstTransactionTime, secondTransactionTime, thirdTransactionTime}
			expectedTransactionNatures := []model.TransactionNature{"credit", "debit", "credit"}
			expectedTransactionTypes := []model.TransactionType{"realized", "realized", "realized"}
			expectedResults := createExpectedAccountTranResult(
				tomAccountName,
				samAccountName,
				expectedCreatedAts,
				tranAmounts,
				expectedTransactionNatures,
				expectedTransactionTypes,
			)
			assertExpectedMatchesResult(t, expectedResults, res)
		},
	)

	t.Run("Does not return transactions outside the specified time range", func(t *testing.T) {
		setupGetAccountTransactionsTestCase(accCollection, tranCollection, ctx, t)
		_, err := tranCollection.InsertMany(ctx, transactionsInput)
		if err != nil {
			t.Errorf("Error inserting transactions: %v", err)
		}

		accountService := setupAccountService(mongoClient, tranCollection, accCollection)

		input := model.TransactionsForBankAccountInput{
			BankAccountId: tomAccountName,
			FromTime:      secondTransactionTime,
			ToTime:        thirdTransactionTime,
		}
		res, _ := accountService.GetBankAccountTransactions(&input, ctx)
		expectedCreatedAts := []time.Time{secondTransactionTime, thirdTransactionTime}
		expectedTransactionNatures := []model.TransactionNature{"debit", "credit"}
		expectedTransactionTypes := []model.TransactionType{"realized", "realized"}
		expectedResults := createExpectedAccountTranResult(
			tomAccountName,
			samAccountName,
			expectedCreatedAts,
			tranAmounts[1:],
			expectedTransactionNatures,
			expectedTransactionTypes,
		)
		assertExpectedMatchesResult(t, expectedResults, res)
	})
}

func setupGetAccountTransactionsTestCase(
	accCollection *mongo.Collection,
	tranCollection *mongo.Collection,
	ctx context.Context,
	t *testing.T,
) {
	utils.CleanupCollection(tranCollection, ctx)
	utils.CleanupCollection(accCollection, ctx)

	_, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	_, samErr := accCollection.InsertOne(ctx, utils.SamAccountDetails)
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}
}

var firstTransactionTime = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
var secondTransactionTime = time.Date(2021, time.January, 10, 0, 0, 0, 0, time.UTC)
var thirdTransactionTime = time.Date(2021, time.January, 20, 0, 0, 0, 0, time.UTC)

func makeTransactionsInput(
	tomAccountId primitive.ObjectID,
	samAccountId primitive.ObjectID,
	tranAmounts []primitive.Decimal128,
) []interface{} {
	return []interface{}{
		mongodb.MongoTransactionInput{
			FromBankAccountId: tomAccountId,
			ToBankAccountId:   samAccountId,
			Amount:            tranAmounts[0],
			Type:              "realized",
			CreatedAt:         pkgutils.TimeToTimestamp(firstTransactionTime),
		},
		mongodb.MongoTransactionInput{
			FromBankAccountId: samAccountId,
			ToBankAccountId:   tomAccountId,
			Amount:            tranAmounts[1],
			Type:              "realized",
			CreatedAt:         pkgutils.TimeToTimestamp(secondTransactionTime),
		},
		mongodb.MongoTransactionInput{
			FromBankAccountId: tomAccountId,
			ToBankAccountId:   samAccountId,
			Amount:            tranAmounts[2],
			Type:              "realized",
			CreatedAt:         pkgutils.TimeToTimestamp(thirdTransactionTime),
		},
	}
}

func createExpectedAccountTranResult(
	accountId string,
	otherAccountId string,
	expectedCreatedAts []time.Time,
	tranAmounts []decimal.Decimal,
	transactionNatures []model.TransactionNature,
	transactionTypes []model.TransactionType,
) []model.BankAccountTransactionOutput {
	expectedResults := make([]model.BankAccountTransactionOutput, len(tranAmounts))
	for i, _ := range tranAmounts {
		expectedResults[i] = model.BankAccountTransactionOutput{
			BankAccountId:      accountId,
			TransactionNature:  transactionNatures[i],
			TransactionType:    transactionTypes[i],
			OtherBankAccountId: otherAccountId,
			Amount:             tranAmounts[i],
			CreatedAt:          expectedCreatedAts[i],
		}
	}
	return expectedResults
}

func assertExpectedMatchesResult(
	t *testing.T,
	expectedResults []model.BankAccountTransactionOutput,
	res []model.BankAccountTransactionOutput,
) {
	for i, _ := range expectedResults {
		assert.Equal(t, expectedResults[i].BankAccountId, res[i].BankAccountId)
		assert.Equal(t, expectedResults[i].TransactionType, res[i].TransactionType)
		assert.Equal(t, expectedResults[i].OtherBankAccountId, res[i].OtherBankAccountId)
		assert.Equal(t, expectedResults[i].Amount.String(), res[i].Amount.String())
		assert.Equal(t, expectedResults[i].CreatedAt.Unix(), res[i].CreatedAt.Unix())
	}
}
