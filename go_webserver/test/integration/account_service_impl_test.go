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

func GetAccountDetails(t *testing.T) {
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
