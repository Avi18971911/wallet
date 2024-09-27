package integration

import (
	"context"
	"github.com/shopspring/decimal"
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
	utils.CleanupMigrations(tranCollection, ctx)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection("account")
	utils.CleanupMigrations(accCollection, ctx)
	_, tomErr := accCollection.InsertOne(ctx, utils.TomAccountDetails)
	if tomErr != nil {
		t.Errorf("Error inserting Tom's record %v", tomErr)
	}
	_, samErr := accCollection.InsertOne(ctx, utils.SamAccountDetails)
	if samErr != nil {
		t.Errorf("Error inserting Tom's record %v", samErr)
	}

	tomAccountName, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.Accounts[0].Id)
	samAccountName, _ := pkgutils.ObjectIdToString(utils.SamAccountDetails.Accounts[0].Id)
	service := setupTransactionService(mongoClient, tranCollection, accCollection)
	transferAmount := decimal.NewFromFloatWithExponent(1000.0, -2)
	input := model.TransactionDetails{
		ToAccount:   samAccountName,
		FromAccount: tomAccountName,
		Amount:      transferAmount,
	}

	t.Run("Should be able to insert transactions", func(t *testing.T) {
		err := service.AddTransaction(input.ToAccount, input.FromAccount, input.Amount.String(), ctx)
		assert.Nil(t, err)
		samFind, tomFind := mongodb.MongoAccountOutput{}, mongodb.MongoAccountOutput{}
		err = accCollection.FindOne(
			ctx, bson.M{"accounts._id": utils.SamAccountDetails.Accounts[0].Id},
		).Decode(&samFind)
		if err != nil {
			t.Errorf("Error in finding Sam's account details: %v", err)
		}
		err = accCollection.FindOne(
			ctx, bson.M{"accounts._id": utils.TomAccountDetails.Accounts[0].Id},
		).Decode(&tomFind)
		if err != nil {
			t.Errorf("Error in finding Tom's's account details: %v", err)
		}
		samBalance, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(
			utils.SamAccountDetails.Accounts[0].AvailableBalance,
		)
		samBalance = samBalance.Add(transferAmount)
		tomBalance, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(
			utils.TomAccountDetails.Accounts[0].AvailableBalance,
		)
		tomBalance = tomBalance.Sub(transferAmount)
		assert.Equal(
			t,
			samBalance.String(),
			samFind.Accounts[0].AvailableBalance.String(),
		)
		assert.Equal(
			t,
			tomBalance.String(),
			tomFind.Accounts[0].AvailableBalance.String(),
		)

		var tranRes = mongodb.MongoTransactionInput{}
		err = tranCollection.FindOne(
			ctx, bson.M{"fromAccount": utils.TomAccountDetails.Accounts[0].Id},
		).Decode(&tranRes)
		if err != nil {
			t.Errorf("Error in decoding transaction result into tranRes: %v", err)
		}
		assert.Equal(t, transferAmount.String(), tranRes.Amount.String())
	})

	t.Run("Should not be able to insert transactions with insufficient balance", func(t *testing.T) {
		reallyHighAmount := "99999999.99"
		err := service.AddTransaction(input.ToAccount, input.FromAccount, reallyHighAmount, ctx)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "insufficient balance")
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
