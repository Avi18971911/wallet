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
	tomAccountId := primitive.NewObjectID()
	samAccountId := primitive.NewObjectID()
	accCollection.InsertOne(ctx, bson.M{"availableBalance": 1000.0, "_id": tomAccountId})
	accCollection.InsertOne(ctx, bson.M{"availableBalance": 1000.0, "_id": samAccountId})

	tomAccountName, err := pkgutils.ObjectIdToString(tomAccountId)
	samAccountName, err := pkgutils.ObjectIdToString(samAccountId)

	tranId1 := primitive.NewObjectID()
	tranId2 := primitive.NewObjectID()
	tranId3 := primitive.NewObjectID()
	tranString1, err := pkgutils.ObjectIdToString(tranId1)
	tranString2, err := pkgutils.ObjectIdToString(tranId2)
	tranString3, err := pkgutils.ObjectIdToString(tranId3)

	_, err = tranCollection.InsertMany(ctx, bson.A{
		bson.M{"fromAccount": tomAccountId, "toAccount": samAccountId, "amount": 50.32, "_id": tranId1},
		bson.M{"fromAccount": samAccountId, "toAccount": tomAccountId, "amount": 23.89, "_id": tranId2},
		bson.M{"fromAccount": tomAccountId, "toAccount": samAccountId, "amount": 10.88, "_id": tranId3},
	})
	if err != nil {
		t.Errorf("Error inserting transactions: %v", err)
	}

	accountService := setupAccountService(mongoClient, tranCollection, accCollection)

	res, err := accountService.GetAccountTransactions(tomAccountName, ctx)
	expectedCreatedAt := res[0].CreatedAt
	expectedResults := []model.AccountTransaction{
		{
			Id:              tranString1,
			AccountId:       tomAccountName,
			TransactionType: "debit",
			OtherAccountId:  samAccountName,
			Amount:          50.32,
			CreatedAt:       expectedCreatedAt,
		},
		{
			Id:              tranString2,
			AccountId:       tomAccountName,
			TransactionType: "credit",
			OtherAccountId:  samAccountName,
			Amount:          23.89,
			CreatedAt:       expectedCreatedAt,
		},
		{
			Id:              tranString3,
			AccountId:       tomAccountName,
			TransactionType: "debit",
			OtherAccountId:  samAccountName,
			Amount:          10.88,
			CreatedAt:       expectedCreatedAt,
		},
	}
	assert.ElementsMatch(t, expectedResults, res)
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
