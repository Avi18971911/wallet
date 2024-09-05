package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/transactional"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/test/utils"
)

func TestGetAccountDetails(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection("transaction")
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")

	t.Run("Allows the retrieval of account details from an inserted account record", func(t *testing.T) {
		tomRes, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
		if tomErr != nil {
			t.Errorf("Error inserting Tom's record %v", tomErr)
		}
		tomAccountId, _ := pkgutils.ObjectIdToString(tomRes.InsertedID)
		service := setupAccountService(mongoClient, tranCollection, accCollection)

		accountDetails, err := service.GetAccountDetails(tomAccountId, ctx)
		if err != nil {
			t.Errorf("Error getting Tom's accountDetails: %v", err)
		}
		assert.Equal(t, 1000.0, accountDetails.AvailableBalance)
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(t, pkgutils.TimestampToTime(utils.TomAccountDetails.CreatedAt), accountDetails.CreatedAt)
		assert.Equal(t, utils.TomAccountDetails.Password, accountDetails.Password)
		assert.Equal(t, tomAccountId, accountDetails.Id)
	})
}

func TestGetAccountTransactions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection("transaction")
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")
	tomRes, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	samRes, samErr := accCollection.InsertOne(ctx, utils.SamAccountDetails)
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}

	tomAccountName, _ := pkgutils.ObjectIdToString(tomRes.InsertedID)
	samAccountName, _ := pkgutils.ObjectIdToString(samRes.InsertedID)
	tomObjectId, _ := pkgutils.StringToObjectId(tomAccountName)
	samObjectId, _ := pkgutils.StringToObjectId(samAccountName)

	tranAmounts := []float64{50.32, 23.89, 10.88}
	tranIds := make([]primitive.ObjectID, 3)
	tranStrings := make([]string, 3)
	for i, _ := range tranIds {
		tranIds[i] = primitive.NewObjectID()
		tranStrings[i], _ = pkgutils.ObjectIdToString(tranIds[i])
	}
	transactionsInput := makeTransactionsInput(tomObjectId, samObjectId, tranAmounts, tranIds)

	t.Run(
		"Allows the insertion of transactions and the retrieval of all transactions from an account",
		func(t *testing.T) {
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
			assert.ElementsMatch(t, expectedResults, res)
		},
	)
}

func TestLogins(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection("transaction")
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")

	_, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	_, samErr := accCollection.InsertOne(ctx, utils.SamAccountDetails)
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}

	t.Run("Allows the login of a user with the correct password", func(t *testing.T) {
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		accountDetails, err := service.Login(utils.TomAccountDetails.Username, utils.TomAccountDetails.Password, ctx)
		if err != nil {
			t.Fatalf("Error logging in: %v", err)
		}
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(t, utils.TomAccountDetails.StartingBalance, accountDetails.AvailableBalance)
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

func makeTransactionsInput(
	tomAccountId primitive.ObjectID,
	samAccountId primitive.ObjectID,
	tranAmounts []float64,
	tranIds []primitive.ObjectID,
) bson.A {
	return bson.A{
		bson.M{
			"fromAccount": tomAccountId,
			"toAccount":   samAccountId,
			"amount":      tranAmounts[0],
			"_id":         tranIds[0],
			"_createdAt":  pkgutils.GetCurrentTimestamp(),
		},
		bson.M{
			"fromAccount": samAccountId,
			"toAccount":   tomAccountId,
			"amount":      tranAmounts[1],
			"_id":         tranIds[1],
			"_createdAt":  pkgutils.GetCurrentTimestamp(),
		},
		bson.M{
			"fromAccount": tomAccountId,
			"toAccount":   samAccountId,
			"amount":      tranAmounts[2],
			"_id":         tranIds[2],
			"_createdAt":  pkgutils.GetCurrentTimestamp(),
		},
	}
}

func createExpectedAccountTranResult(
	accountId string,
	otherAccountId string,
	expectedCreatedAt time.Time,
	tranStrings []string,
	tranAmounts []float64,
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
