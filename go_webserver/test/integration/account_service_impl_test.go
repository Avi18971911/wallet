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
		tomAccountId, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.Accounts[0].Id)
		knownAccountId, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.KnownAccounts[0].Id)
		service := setupAccountService(mongoClient, tranCollection, accCollection)

		accountDetails, err := service.GetAccountDetails(tomAccountId, ctx)
		if err != nil {
			t.Errorf("Error getting Tom's accountDetails: %v", err)
		}
		assert.Equal(
			t,
			utils.TomAccountDetails.Accounts[0].AvailableBalance.String(),
			accountDetails.Accounts[0].AvailableBalance.String(),
		)
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(t, pkgutils.TimestampToTime(utils.TomAccountDetails.CreatedAt), accountDetails.CreatedAt)
		assert.Equal(t, utils.TomAccountDetails.Password, accountDetails.Password)
		assert.Equal(t, knownAccountId, accountDetails.KnownAccounts[0].Id)
		assert.Equal(t, tomAccountId, accountDetails.Accounts[0].Id)
	})
}

func setupGetAccountDetailsTestCase(
	accCollection *mongo.Collection,
	tranCollection *mongo.Collection,
	ctx context.Context,
) {
	utils.CleanupMigrations(tranCollection, ctx)
	utils.CleanupMigrations(accCollection, ctx)
}

func TestGetAccountTransactions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	tomAccountName, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.Accounts[0].Id)
	samAccountName, _ := pkgutils.ObjectIdToString(utils.SamAccountDetails.Accounts[0].Id)
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

	tranIds := make([]primitive.ObjectID, 3)
	tranStrings := make([]string, 3)
	for i, _ := range tranIds {
		tranIds[i] = primitive.NewObjectID()
		tranStrings[i], _ = pkgutils.ObjectIdToString(tranIds[i])
	}
	transactionsInput := makeTransactionsInput(tomObjectId, samObjectId, tranAmountsDecimal128, tranIds)

	t.Run(
		"Allows the insertion of transactions and the retrieval of all transactions from an account",
		func(t *testing.T) {
			setupGetAccountTransactionsTestCase(accCollection, tranCollection, ctx, t)
			_, err := tranCollection.InsertMany(ctx, transactionsInput)
			if err != nil {
				t.Errorf("Error inserting transactions: %v", err)
			}

			accountService := setupAccountService(mongoClient, tranCollection, accCollection)

			res, _ := accountService.GetAccountTransactions(tomAccountName, ctx)
			expectedCreatedAt := res[0].CreatedAt
			expectedTransactionTypes := []string{"credit", "debit", "credit"}
			expectedResults := createExpectedAccountTranResult(
				tomAccountName,
				samAccountName,
				expectedCreatedAt,
				tranStrings,
				tranAmounts,
				expectedTransactionTypes,
			)
			assertExpectedMatchesResult(t, expectedResults, res)
		},
	)
}

func setupGetAccountTransactionsTestCase(
	accCollection *mongo.Collection,
	tranCollection *mongo.Collection,
	ctx context.Context,
	t *testing.T,
) {
	utils.CleanupMigrations(tranCollection, ctx)
	utils.CleanupMigrations(accCollection, ctx)

	_, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	_, samErr := accCollection.InsertOne(ctx, utils.SamAccountDetails)
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}
}

func makeTransactionsInput(
	tomAccountId primitive.ObjectID,
	samAccountId primitive.ObjectID,
	tranAmounts []primitive.Decimal128,
	tranIds []primitive.ObjectID,
) []interface{} {
	return []interface{}{
		mongodb.MongoTransactionInput{
			FromAccount: tomAccountId,
			ToAccount:   samAccountId,
			Amount:      tranAmounts[0],
			Id:          tranIds[0],
			CreatedAt:   pkgutils.GetCurrentTimestamp(),
		},
		mongodb.MongoTransactionInput{
			FromAccount: samAccountId,
			ToAccount:   tomAccountId,
			Amount:      tranAmounts[1],
			Id:          tranIds[1],
			CreatedAt:   pkgutils.GetCurrentTimestamp(),
		},
		mongodb.MongoTransactionInput{
			FromAccount: tomAccountId,
			ToAccount:   samAccountId,
			Amount:      tranAmounts[2],
			Id:          tranIds[2],
			CreatedAt:   pkgutils.GetCurrentTimestamp(),
		},
	}
}

func createExpectedAccountTranResult(
	accountId string,
	otherAccountId string,
	expectedCreatedAt time.Time,
	tranStrings []string,
	tranAmounts []decimal.Decimal,
	transactionTypes []string,
) []model.AccountTransaction {
	expectedResults := make([]model.AccountTransaction, len(tranAmounts))
	for i, _ := range tranAmounts {
		expectedResults[i] = model.AccountTransaction{
			Id:              tranStrings[i],
			AccountId:       accountId,
			TransactionType: transactionTypes[i],
			OtherAccountId:  otherAccountId,
			Amount:          tranAmounts[i],
			CreatedAt:       expectedCreatedAt,
		}
	}
	return expectedResults
}

func assertExpectedMatchesResult(
	t *testing.T,
	expectedResults []model.AccountTransaction,
	res []model.AccountTransaction,
) {
	for i, _ := range expectedResults {
		assert.Equal(t, expectedResults[i].Id, res[i].Id)
		assert.Equal(t, expectedResults[i].AccountId, res[i].AccountId)
		assert.Equal(t, expectedResults[i].TransactionType, res[i].TransactionType)
		assert.Equal(t, expectedResults[i].OtherAccountId, res[i].OtherAccountId)
		assert.Equal(t, expectedResults[i].Amount.String(), res[i].Amount.String())
		assert.Equal(t, expectedResults[i].CreatedAt, res[i].CreatedAt)
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
			t, utils.TomAccountDetails.Accounts[0].AvailableBalance.String(),
			accountDetails.Accounts[0].AvailableBalance.String(),
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
	utils.CleanupMigrations(tranCollection, ctx)
	utils.CleanupMigrations(accCollection, ctx)

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
