package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
	"webserver/migrations/versions/schema"
	"webserver/test/utils"
)

func TestLogins(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	t.Run("Allows the login of a user with the correct password", func(t *testing.T) {
		setupLoginTestCase(accCollection, tranCollection, ctx, t)
		service := setupAccountService(mongoClient, tranCollection, accCollection)
		accountDetails, err := service.Login(utils.TomAccountDetails.Username, utils.TomAccountDetails.Password, ctx)
		if err != nil {
			t.Fatalf("Error logging in: %v", err)
		}
		assert.Equal(t, utils.TomAccountDetails.Username, accountDetails.Username)
		assert.Equal(
			t, utils.TomAccountDetails.BankAccounts[0].AvailableBalance.String(),
			accountDetails.BankAccounts[0].AvailableBalance.String(),
		)
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

func setupLoginTestCase(
	accCollection *mongo.Collection,
	tranCollection *mongo.Collection,
	ctx context.Context,
	t *testing.T,
) {
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
