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
	"webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/infrastructure/transactional"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/migrations/versions/schema"
	"webserver/test/utils"
)

func TestGetAccountDetails(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	t.Run("Allows the retrieval of account details from an inserted account record", func(t *testing.T) {
		setupGetAccountDetailsTestCase(accCollection, tranCollection, ctx)
		_, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
		if tomErr != nil {
			t.Errorf("Error inserting Tom's record %v", tomErr)
		}
		tomAccountId, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.BankAccounts[0].Id)
		knownAccountId, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.KnownBankAccounts[0].Id)
		service := setupAccountService(mongoClient, tranCollection, accCollection)

		accountDetails, err := service.GetAccountDetailsFromBankAccountId(tomAccountId, ctx)
		if err != nil {
			t.Errorf("Error getting Tom's accountDetails: %v", err)
		}
		assert.Equal(
			t,
			utils.TomAccountDetails.BankAccounts[0].AvailableBalance.String(),
			accountDetails.BankAccounts[0].AvailableBalance.String(),
		)
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(t, pkgutils.TimestampToTime(utils.TomAccountDetails.CreatedAt), accountDetails.CreatedAt)
		assert.Equal(t, utils.TomAccountDetails.Password, accountDetails.Password)
		assert.Equal(t, knownAccountId, accountDetails.KnownBankAccounts[0].Id)
		assert.Equal(t, tomAccountId, accountDetails.BankAccounts[0].Id)
	})
}

func setupGetAccountDetailsTestCase(
	accCollection *mongo.Collection,
	tranCollection *mongo.Collection,
	ctx context.Context,
) {
	utils.CleanupCollection(tranCollection, ctx)
	utils.CleanupCollection(accCollection, ctx)
}

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

func TestLogins(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	t.Run("Allows the login of a user with the correct password", func(t *testing.T) {
		setupLoginTestCase(accCollection, tranCollection, ctx, t)
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		accountDetails, err := service.Login(utils.TomAccountDetails.Username, utils.TomAccountDetails.Password, ctx)
		if err != nil {
			t.Fatalf("Error logging in: %v", err)
		}
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(
			t, utils.TomAccountDetails.BankAccounts[0].AvailableBalance.String(),
			accountDetails.BankAccounts[0].AvailableBalance.String(),
		)
	})

	t.Run("Does not allow the login of a user with the incorrect password", func(t *testing.T) {
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		_, err := service.Login(utils.TomAccountDetails.Username, "wrongpassword", ctx)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid username or password")
	})

	t.Run("Does not allow the login of a user with the incorrect username", func(t *testing.T) {
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		_, err := service.Login("wrongusername", utils.TomAccountDetails.Password, ctx)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid username or password")
	})

}

func setupLoginTestCase(
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

func setupAccountService(
	mongoClient *mongo.Client,
	tranCollection *mongo.Collection,
	accCollection *mongo.Collection,
) *services.AccountServiceImpl {
	tr := repositories.CreateNewTransactionRepositoryMongodb(tranCollection)
	ar := repositories.CreateNewAccountRepositoryMongodb(accCollection)
	tran := transactional.NewMongoTransactional(mongoClient)
	service := services.CreateNewAccountServiceImpl(ar, tr, tran)
	return service
}
