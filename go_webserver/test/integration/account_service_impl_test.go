package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
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
	tomRes, tomErr := accCollection.InsertOne(ctx, bson.M{"availableBalance": 1000.0, "accountHolderFirstName": "Tom"})
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
	tomRes, tomErr := accCollection.InsertOne(ctx, bson.M{"availableBalance": 1000.0, "accountHolderFirstName": "Tom"})
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	samRes, samErr := accCollection.InsertOne(ctx, bson.M{"availableBalance": 1000.0, "accountHolderFirstName": "Sam"})
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}

	tomAccountName, _ := pkgutils.ObjectIdToString(tomRes.InsertedID)
	// samAccountName, _ := pkgutils.ObjectIdToString(samRes.InsertedID)

	_, err := tranCollection.InsertMany(ctx, bson.A{
		bson.M{"fromAccount": tomRes.InsertedID, "toAccount": samRes.InsertedID, "amount": 50.32},
		bson.M{"fromAccount": samRes.InsertedID, "toAccount": tomRes.InsertedID, "amount": 23.89},
		bson.M{"fromAccount": tomRes.InsertedID, "toAccount": samRes.InsertedID, "amount": 10.88},
	})
	if err != nil {
		t.Errorf("Error inserting transactions: %v", err)
	}

	accountService := setupAccountService(mongoClient, tranCollection, accCollection)

	res, err := accountService.GetAccountTransactions(tomAccountName, ctx)
	assert.Nil(t, err)
	assert.Equal(t, 50.32, res[0].Amount)
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
