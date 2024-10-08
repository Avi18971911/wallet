package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
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
		tomAccountId, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.BankAccounts[0].Id)
		knownAccountId, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.KnownBankAccounts[0].Id)
		service := setupAccountService(mongoClient, tranCollection, accCollection)

		accountDetails, err := service.GetAccountDetailsFromBankAccountId(tomAccountId, ctx)
		if err != nil {
			t.Errorf("Error getting Tom's accountDetails: %v", err)
		}
		assert.Equal(
			t,
			utils.TomAccountDetails.BankAccounts[0].AvailableBalance.String(),
			accountDetails.BankAccounts[0].AvailableBalance.String(),
		)
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(t, pkgutils.TimestampToTime(utils.TomAccountDetails.CreatedAt), accountDetails.CreatedAt)
		assert.Equal(t, utils.TomAccountDetails.Password, accountDetails.Password)
		assert.Equal(t, knownAccountId, accountDetails.KnownBankAccounts[0].Id)
		assert.Equal(t, tomAccountId, accountDetails.BankAccounts[0].Id)
	})
}

func setupGetAccountDetailsTestCase(
	accCollection *mongo.Collection,
	tranCollection *mongo.Collection,
	ctx context.Context,
) {
	utils.CleanupCollection(tranCollection, ctx)
	utils.CleanupCollection(accCollection, ctx)
}
