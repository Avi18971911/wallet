package integration

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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

func TestAddTransaction(t *testing.T) {
	if mongoClient == nil {
		t.Error("mongoClient is uninitialized or otherwise nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	tomAccountName, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.BankAccounts[0].Id)
	samAccountName, _ := pkgutils.ObjectIdToString(utils.SamAccountDetails.BankAccounts[0].Id)

	service := setupTransactionService(mongoClient, tranCollection, accCollection)
	transferAmount, _ := decimal.NewFromString("100.00")
	input := model.TransactionDetailsInput{
		ToBankAccountId:   samAccountName,
		FromBankAccountId: tomAccountName,
		Amount:            transferAmount,
		Type:              model.Realized,
	}

	t.Run("Should be able to insert transactions", func(t *testing.T) {
		setupAddTransactionTestCase(tranCollection, accCollection, ctx, t)
		err := service.AddTransaction(input, ctx)
		assert.Nil(t, err)
		samFind, tomFind := mongodb.MongoAccountOutput{}, mongodb.MongoAccountOutput{}
		err = accCollection.FindOne(
			ctx, bson.M{"bankAccounts._id": utils.SamAccountDetails.BankAccounts[0].Id},
		).Decode(&samFind)
		if err != nil {
			t.Errorf("Error in finding Sam's account details: %v", err)
		}
		err = accCollection.FindOne(
			ctx, bson.M{"bankAccounts._id": utils.TomAccountDetails.BankAccounts[0].Id},
		).Decode(&tomFind)
		if err != nil {
			t.Errorf("Error in finding Tom's's account details: %v", err)
		}
		samBalance, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(
			utils.SamAccountDetails.BankAccounts[0].AvailableBalance,
		)
		samBalance = samBalance.Add(transferAmount)
		tomBalance, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(
			utils.TomAccountDetails.BankAccounts[0].AvailableBalance,
		)
		tomBalance = tomBalance.Sub(transferAmount)
		assert.Equal(
			t,
			samBalance.String(),
			samFind.BankAccounts[0].AvailableBalance.String(),
		)
		assert.Equal(
			t,
			tomBalance.String(),
			tomFind.BankAccounts[0].AvailableBalance.String(),
		)
		assert.Equal(
			t,
			samBalance.String(),
			samFind.BankAccounts[0].PendingBalance.String(),
		)
		assert.Equal(
			t,
			tomBalance.String(),
			tomFind.BankAccounts[0].PendingBalance.String(),
		)
		var tranRes = mongodb.MongoTransactionInput{}
		err = tranCollection.FindOne(
			ctx, bson.M{"fromBankAccountId": utils.TomAccountDetails.BankAccounts[0].Id},
		).Decode(&tranRes)
		if err != nil {
			t.Errorf("Error in decoding transaction result into tranRes: %v", err)
		}
		assert.Equal(t, transferAmount.String(), tranRes.Amount.String())
	})

	t.Run("Should not be able to insert transactions with insufficient balance", func(t *testing.T) {
		setupAddTransactionTestCase(tranCollection, accCollection, ctx, t)
		reallyHighAmount := "99999999.99"
		reallyHighInput := model.TransactionDetailsInput{
			ToBankAccountId:   samAccountName,
			FromBankAccountId: tomAccountName,
			Amount:            decimal.RequireFromString(reallyHighAmount),
			Type:              model.Realized,
		}
		err := service.AddTransaction(reallyHighInput, ctx)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "insufficient balance in BankAccount "+tomAccountName)
	})

	t.Run("Should not carry out transaction if there is an error", func(t *testing.T) {
		setupAddTransactionTestCase(tranCollection, accCollection, ctx, t)
		reallyHighAmount := "99999999.99"
		reallyHighInput := model.TransactionDetailsInput{
			ToBankAccountId:   samAccountName,
			FromBankAccountId: tomAccountName,
			Amount:            decimal.RequireFromString(reallyHighAmount),
		}
		err := service.AddTransaction(reallyHighInput, ctx)
		assert.NotNil(t, err)
		samDetails, tomDetails := mongodb.MongoAccountOutput{}, mongodb.MongoAccountOutput{}

		err = accCollection.FindOne(
			ctx, bson.M{"bankAccounts._id": utils.SamAccountDetails.BankAccounts[0].Id},
		).Decode(&samDetails)

		if err != nil {
			t.Errorf("Error in finding Sam's account details: %v", err)
		}

		err = accCollection.FindOne(
			ctx, bson.M{"bankAccounts._id": utils.TomAccountDetails.BankAccounts[0].Id},
		).Decode(&tomDetails)

		if err != nil {
			t.Errorf("Error in finding Tom's's account details: %v", err)
		}

		assert.Equal(
			t,
			utils.SamAccountDetails.BankAccounts[0].AvailableBalance.String(),
			samDetails.BankAccounts[0].AvailableBalance.String(),
		)
		assert.Equal(
			t,
			utils.TomAccountDetails.BankAccounts[0].AvailableBalance.String(),
			tomDetails.BankAccounts[0].AvailableBalance.String(),
		)
	})

	t.Run("Pending transactions should not be realized", func(t *testing.T) {
		setupAddTransactionTestCase(tranCollection, accCollection, ctx, t)
		pendingInput := model.TransactionDetailsInput{
			ToBankAccountId:   samAccountName,
			FromBankAccountId: tomAccountName,
			Amount:            transferAmount,
			Type:              model.Pending,
		}
		err := service.AddTransaction(pendingInput, ctx)
		assert.Nil(t, err)
		samFind, tomFind := mongodb.MongoAccountOutput{}, mongodb.MongoAccountOutput{}
		err = accCollection.FindOne(
			ctx, bson.M{"bankAccounts._id": utils.SamAccountDetails.BankAccounts[0].Id},
		).Decode(&samFind)
		if err != nil {
			t.Errorf("Error in finding Sam's account details: %v", err)
		}
		err = accCollection.FindOne(
			ctx, bson.M{"bankAccounts._id": utils.TomAccountDetails.BankAccounts[0].Id},
		).Decode(&tomFind)
		if err != nil {
			t.Errorf("Error in finding Tom's's account details: %v", err)
		}
		samBalance, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(
			utils.SamAccountDetails.BankAccounts[0].PendingBalance,
		)
		samBalance = samBalance.Add(transferAmount)
		tomBalance, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(
			utils.TomAccountDetails.BankAccounts[0].PendingBalance,
		)
		tomBalance = tomBalance.Sub(transferAmount)
		assert.Equal(
			t,
			samBalance.String(),
			formatDecimalString(samFind.BankAccounts[0].PendingBalance.String()),
		)
		assert.Equal(
			t,
			tomBalance.String(),
			formatDecimalString(tomFind.BankAccounts[0].PendingBalance.String()),
		)
		assert.NotEqual(
			t,
			samFind.BankAccounts[0].AvailableBalance.String(),
			samFind.BankAccounts[0].PendingBalance.String(),
		)
		assert.NotEqual(
			t,
			tomFind.BankAccounts[0].AvailableBalance.String(),
			tomFind.BankAccounts[0].PendingBalance.String(),
		)
		var tranRes = mongodb.MongoTransactionInput{}
		err = tranCollection.FindOne(
			ctx, bson.M{"fromBankAccountId": utils.TomAccountDetails.BankAccounts[0].Id},
		).Decode(&tranRes)
		if err != nil {
			t.Errorf("Error in decoding transaction result into tranRes: %v", err)
		}
		assert.Equal(t, string(model.Pending), tranRes.Type)
	})
}

func setupAddTransactionTestCase(
	tranCollection *mongo.Collection,
	accCollection *mongo.Collection,
	ctx context.Context,
	t *testing.T,
) {
	log.Printf("Cleaning up the database")
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

func formatDecimalString(input string) string {
	return decimal.RequireFromString(input).String()
}
