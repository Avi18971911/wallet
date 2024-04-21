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
	tomRes, tomErr := accCollection.InsertOne(ctx, bson.M{"availableBalance": 1000.0})
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
	assert.Equal(t, tomAccountName, accountDetails.Id)
}

func TestGetAccountTransactions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection("transaction")
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")
	baseAmounts := []float64{1000.0, 1000.0}
	accountIds, err := createAccounts(accCollection, ctx, baseAmounts)
	tomAccountId, samAccountId := accountIds[0], accountIds[1]

	tomAccountName, err := pkgutils.ObjectIdToString(tomAccountId)
	samAccountName, err := pkgutils.ObjectIdToString(samAccountId)

	tranAmounts := []float64{50.32, 23.89, 10.88}

	tranId1 := primitive.NewObjectID()
	tranId2 := primitive.NewObjectID()
	tranId3 := primitive.NewObjectID()
	tranString1, err := pkgutils.ObjectIdToString(tranId1)
	tranString2, err := pkgutils.ObjectIdToString(tranId2)
	tranString3, err := pkgutils.ObjectIdToString(tranId3)
	tranStrings := []string{tranString1, tranString2, tranString3}

	_, err = tranCollection.InsertMany(ctx, bson.A{
		bson.M{"fromAccount": tomAccountId, "toAccount": samAccountId, "amount": tranAmounts[0], "_id": tranId1},
		bson.M{"fromAccount": samAccountId, "toAccount": tomAccountId, "amount": tranAmounts[1], "_id": tranId2},
		bson.M{"fromAccount": tomAccountId, "toAccount": samAccountId, "amount": tranAmounts[2], "_id": tranId3},
	})
	if err != nil {
		t.Errorf("Error inserting transactions: %v", err)
	}

	accountService := setupAccountService(mongoClient, tranCollection, accCollection)

	res, err := accountService.GetAccountTransactions(tomAccountName, ctx)
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
) (accountIds []primitive.ObjectID, err error) {
	res := make([]primitive.ObjectID, len(amounts))
	for i, _ := range amounts {
		currentAmount := amounts[i]
		res[i] = primitive.NewObjectID()
		_, err := insertAccount(accCollection, ctx, res[i], currentAmount)
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
) (*mongo.InsertOneResult, error) {
	return accCollection.InsertOne(ctx, bson.M{"availableBalance": amount, "_id": accountId})
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
