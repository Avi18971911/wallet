//go:build test

package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/repositories"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/infrastructure/transactional"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/test/utils"
)

func TestAddTransaction(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection("transaction")
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")
	tomRes, tomErr := accCollection.InsertOne(ctx, bson.M{
		"availableBalance": 1000.0,
		"username":         "Tom",
		"password":         "pass",
		"_createdAt":       pkgutils.GetCurrentTimestamp(),
	})
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	samRes, samErr := accCollection.InsertOne(ctx, bson.M{
		"availableBalance": 1000.0,
		"username":         "Sam",
		"password":         "word",
		"_createdAt":       pkgutils.GetCurrentTimestamp(),
	})
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}

	tomAccountName, _ := pkgutils.ObjectIdToString(tomRes.InsertedID)
	samAccountName, _ := pkgutils.ObjectIdToString(samRes.InsertedID)
	service := setupTransactionService(mongoClient, tranCollection, accCollection)
	transferAmount, baseAmount := 50.42, 1000.0
	input := model.TransactionDetails{
		ToAccount:   samAccountName,
		FromAccount: tomAccountName,
		Amount:      transferAmount,
	}

	t.Run("Should be able to insert transactions with the required fields", func(t *testing.T) {
		err := service.AddTransaction(input.ToAccount, input.FromAccount, input.Amount, ctx)
		assert.Nil(t, err)
		samFind, tomFind := mongodb.MongoAccountDetails{}, mongodb.MongoAccountDetails{}
		err = accCollection.FindOne(ctx, bson.M{"_id": samRes.InsertedID}).Decode(&samFind)
		if err != nil {
			t.Errorf("Error in finding Sam's account details: %v", err)
		}
		err = accCollection.FindOne(ctx, bson.M{"_id": tomRes.InsertedID}).Decode(&tomFind)
		if err != nil {
			t.Errorf("Error in finding Tom's's account details: %v", err)
		}
		assert.Equal(t, baseAmount+transferAmount, samFind.AvailableBalance)
		assert.Equal(t, baseAmount-transferAmount, tomFind.AvailableBalance)

		var tranRes = mongodb.MongoTransactionDetails{}
		err = tranCollection.FindOne(ctx, bson.M{"fromAccount": tomRes.InsertedID}).Decode(&tranRes)
		if err != nil {
			t.Errorf("Error in decoding transaction result into tranRes: %v", err)
		}
		assert.Equal(t, transferAmount, tranRes.Amount)
	})
}

func setupTransactionService(
	mongoClient *mongo.Client,
	tranCollection *mongo.Collection,
	accCollection *mongo.Collection,
) *services.TransactionServiceImpl {
	tr := repositories.CreateNewTransactionRepositoryMongodb(tranCollection)
	ar := repositories.CreateNewAccountRepositoryMongodb(accCollection)
	tran := transactional.NewMongoTransactional(mongoClient)
	service := services.CreateNewTransactionServiceImpl(tr, ar, tran)
	return service
}
