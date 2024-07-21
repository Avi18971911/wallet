//go:build test

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
	username := "user"
	password := "pass"
	ts := pkgutils.GetCurrentTimestamp()

	t.Run("Allows the insertion of an account with the required details", func(t *testing.T) {
		input := bson.M{
			"availableBalance": 1000.0,
			"username":         username,
			"password":         password,
			"_createdAt":       ts,
		}
		tomRes, tomErr := accCollection.InsertOne(ctx, input)
		if tomErr != nil {
			t.Errorf("Error inserting Tom's record %v", tomErr)
		}
		tomAccountName, _ := pkgutils.ObjectIdToString(tomRes.InsertedID)
		service := setupAccountService(mongoClient, tranCollection, accCollection)

		accountDetails, err := service.GetAccountDetails(tomAccountName, ctx)
		if err != nil {
			t.Errorf("Error getting Tom's accountDetails: %v", err)
		}
		assert.Equal(t, 1000.0, accountDetails.AvailableBalance)
		assert.Equal(t, username, accountDetails.Username)
		assert.Equal(t, pkgutils.TimestampToTime(ts), accountDetails.CreatedAt)
		assert.Equal(t, tomAccountName, accountDetails.Id)
	})
}

func TestGetAccountTransactions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection("transaction")
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")
	baseAmounts := []float64{1000.0, 1000.0}
	users := []string{"Tom", "Sam"}
	passwords := []string{"pass", "word"}
	accountIds, err := createAccounts(accCollection, ctx, baseAmounts, users, passwords)
	tomAccountId, samAccountId := accountIds[0], accountIds[1]

	tomAccountName, err := pkgutils.ObjectIdToString(tomAccountId)
	samAccountName, err := pkgutils.ObjectIdToString(samAccountId)

	tranAmounts := []float64{50.32, 23.89, 10.88}
	tranIds := make([]primitive.ObjectID, 3)
	tranStrings := make([]string, 3)
	for i, _ := range tranIds {
		tranIds[i] = primitive.NewObjectID()
		tranStrings[i], _ = pkgutils.ObjectIdToString(tranIds[i])
	}
	transactionsInput := makeTransactionsInput(tomAccountId, samAccountId, tranAmounts, tranIds)

	t.Run(
		"Allows the insertion of transactions and the retrieval of all transactions from an account",
		func(t *testing.T) {
			_, err = tranCollection.InsertMany(ctx, transactionsInput)
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
	baseAmounts := []float64{1000.0, 1000.0}
	users := []string{"Tom", "Sam"}
	passwords := []string{"pass", "word"}
	_, err := createAccounts(accCollection, ctx, baseAmounts, users, passwords)
	if err != nil {
		t.Fatalf("Error creating accounts: %v", err)
	}

	t.Run("Allows the login of a user with the correct password", func(t *testing.T) {
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		accountDetails, err := service.Login(users[0], passwords[0], ctx)
		if err != nil {
			t.Fatalf("Error logging in: %v", err)
		}
		assert.Equal(t, users[0], accountDetails.Username)
		assert.Equal(t, baseAmounts[0], accountDetails.AvailableBalance)
	})

	t.Run("Does not allow the login of a user with the incorrect password", func(t *testing.T) {
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		_, err := service.Login(users[0], "wrongpassword", ctx)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid username or password")
	})

	t.Run("Does not allow the login of a user with the incorrect username", func(t *testing.T) {
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		_, err := service.Login("wrongusername", passwords[0], ctx)
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

func createAccounts(
	accCollection *mongo.Collection,
	ctx context.Context,
	amounts []float64,
	users []string,
	passwords []string,
) (accountIds []primitive.ObjectID, err error) {
	res := make([]primitive.ObjectID, len(amounts))
	for i, _ := range amounts {
		currentAmount := amounts[i]
		res[i] = primitive.NewObjectID()
		_, err := insertAccount(accCollection, ctx, res[i], currentAmount, users[i], passwords[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func insertAccount(
	accCollection *mongo.Collection,
	ctx context.Context,
	accountId primitive.ObjectID,
	amount float64,
	username string,
	password string,
) (*mongo.InsertOneResult, error) {
	return accCollection.InsertOne(ctx, bson.M{
		"availableBalance": amount,
		"_id":              accountId,
		"username":         username,
		"password":         password,
		"_createdAt":       pkgutils.GetCurrentTimestamp(),
	})
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
