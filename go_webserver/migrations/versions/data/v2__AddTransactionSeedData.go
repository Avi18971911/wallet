package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"strconv"
	"time"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
	"webserver/migrations/versions"
	"webserver/migrations/versions/schema"
)

const numTransactions = 100

var transactionIds = make([]primitive.ObjectID, numTransactions)

func createRandomTransactions() []interface{} {
	timeIntervalStart := time.Now().AddDate(-1, 0, 0)
	timeIntervalEnd := time.Now()
	transactions := make([]interface{}, numTransactions)
	for i := 0; i < numTransactions; i++ {
		transactionId := primitive.NewObjectID()
		transactionIds[i] = transactionId
		rndAmountString := strconv.FormatFloat(rand.Float64()*1000, 'f', 2, 64)
		rndAmount, _ := primitive.ParseDecimal128(rndAmountString)
		rndCreatedAt := time.Unix(
			rand.Int63n(timeIntervalEnd.Unix()-timeIntervalStart.Unix())+timeIntervalStart.Unix(),
			0,
		)
		rndBankAccountIndex := rand.Intn(4)
		rndOtherBankAccountIndex := rand.Intn(3)
		var bankAccountId primitive.ObjectID
		var otherBankAccountId primitive.ObjectID
		switch rndBankAccountIndex {
		case 0:
			bankAccountId = ollyAccountId1
			switch rndOtherBankAccountIndex {
			case 0:
				otherBankAccountId = hildaAccountId2
			case 1:
				otherBankAccountId = bobAccountId1
			case 2:
				otherBankAccountId = hildaAccountId1
			}
		case 1:
			bankAccountId = bobAccountId1
			switch rndOtherBankAccountIndex {
			case 0:
				otherBankAccountId = ollyAccountId1
			case 1:
				otherBankAccountId = hildaAccountId2
			case 2:
				otherBankAccountId = hildaAccountId1
			}
		case 2:
			bankAccountId = hildaAccountId1
			switch rndOtherBankAccountIndex {
			case 0, 1:
				otherBankAccountId = ollyAccountId1
			case 2:
				otherBankAccountId = bobAccountId1
			}
		case 3:
			bankAccountId = hildaAccountId2
			switch rndOtherBankAccountIndex {
			case 0, 1:
				otherBankAccountId = bobAccountId1
			case 2:
				otherBankAccountId = ollyAccountId1
			}
		}
		rndTransactionType := rand.Intn(2)
		var transactionType string
		var expiryDate = utils.GetCurrentTimestamp()
		rndTransactionStatus := rand.Intn(3)
		var status string

		if rndTransactionType == 0 {
			transactionType = "realized"
		} else {
			transactionType = "pending"
			switch rndTransactionStatus {
			case 0:
				status = "active"
			case 1:
				status = "applied"
			case 2:
				status = "revoked"
			}
		}

		transactions[i] = mongodb.MongoTransactionOutput{
			Id:                transactionId,
			FromBankAccountId: bankAccountId,
			ToBankAccountId:   otherBankAccountId,
			Type:              transactionType,
			ExpirationDate:    expiryDate,
			Status:            status,
			Amount:            rndAmount,
			CreatedAt:         utils.TimeToTimestamp(rndCreatedAt),
		}
	}

	return transactions
}

var MigrationData2 = versions.Migration{
	Version: "2__Data",
	Up: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		collectionName := schema.TransactionCollectionName
		coll := db.Collection(collectionName)
		transactions := createRandomTransactions()

		_, err := coll.InsertMany(mongoCtx, transactions)
		if err != nil {
			return err
		}

		log.Printf("transaction seed data successfully created")
		return nil
	},
	Down: func(client *mongo.Client, ctx context.Context, databaseName string) error {
		db := client.Database(databaseName)
		mongoCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		collectionName := schema.TransactionCollectionName
		coll := db.Collection(collectionName)

		deleteFilter := bson.M{"_id": bson.M{"$in": transactionIds}}
		_, err := coll.DeleteMany(mongoCtx, deleteFilter)
		if err != nil {
			return err
		}

		return nil
	},
}
