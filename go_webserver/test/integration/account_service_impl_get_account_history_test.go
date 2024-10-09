package integration

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	pkgutils "webserver/internal/pkg/utils"
	"webserver/migrations/versions/schema"
	"webserver/test/utils"
)

func TestGetAccountHistory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	tranCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.TransactionCollectionName)
	accCollection := mongoClient.Database(utils.TestDatabaseName).Collection(schema.AccountCollectionName)

	tomAccountName, _ := pkgutils.ObjectIdToString(utils.TomAccountDetails.BankAccounts[0].Id)
	samAccountName, _ := pkgutils.ObjectIdToString(utils.SamAccountDetails.BankAccounts[0].Id)

	t.Run(
		"Returns the monthly balance history for an account with transactions spanning months",
		func(t *testing.T) {
			setupGetAccountHistoryTestCase(accCollection, tranCollection, ctx, t)
			transactions := createRandomSimpleTransactions(
				tomAccountName,
				samAccountName,
				time.Date(2020, time.December, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.March, 30, 0, 0, 0, 0, time.UTC),
				15,
			)
			_, err := tranCollection.InsertMany(ctx, transactions)
			if err != nil {
				t.Errorf("Error inserting transactions: %v", err)
			}

			accountService := setupAccountService(mongoClient, tranCollection, accCollection)

			input := model.TransactionsForBankAccountInput{
				BankAccountId: tomAccountName,
				ToTime:        time.Date(2021, time.March, 30, 0, 0, 0, 0, time.UTC),
				FromTime:      time.Date(2020, time.December, 1, 0, 0, 0, 0, time.UTC),
			}
			res, _ := accountService.GetAccountHistoryInMonths(&input, ctx)
			monthVals := make([]string, 4)
			for i, month := range res.Months {
				monthVals[i] = month.AvailableBalance.String()
			}
			expectedPendingMarchBalance, _ :=
				pkgutils.FromPrimitiveDecimal128ToDecimal(utils.TomAccountDetails.BankAccounts[0].PendingBalance)
			expectedAvailableMarchBalance, _ :=
				pkgutils.FromPrimitiveDecimal128ToDecimal(utils.TomAccountDetails.BankAccounts[0].AvailableBalance)
			marchTransactions := make([]calculateTransferElement, 0)
			febTransactions := make([]calculateTransferElement, 0)
			janTransactions := make([]calculateTransferElement, 0)
			decemberTransactions := make([]calculateTransferElement, 0)
			for _, transaction := range transactions {
				transaction := transaction.(mongodb.MongoTransactionInput)
				amt, _ := pkgutils.FromPrimitiveDecimal128ToDecimal(transaction.Amount)
				if transaction.ToBankAccountId == utils.TomAccountDetails.BankAccounts[0].Id {
					amt = amt.Neg()
				}
				calculationElement := calculateTransferElement{
					amount:          amt,
					transactionType: model.TransactionType(transaction.Type),
					status:          model.PendingTransactionStatus(transaction.Status),
				}
				switch pkgutils.TimestampToTime(transaction.CreatedAt).Month() {
				case time.December:
					decemberTransactions = append(decemberTransactions, calculationElement)
				case time.January:
					janTransactions = append(janTransactions, calculationElement)
				case time.February:
					febTransactions = append(febTransactions, calculationElement)
				case time.March:
					marchTransactions = append(marchTransactions, calculationElement)
				default:
					t.Errorf("Unexpected month in transaction: %v", transaction)
				}
			}
			expectedPendingFebBalance := expectedPendingMarchBalance
			expectedAvailableFebBalance := expectedAvailableMarchBalance
			expectedAvailableFebBalance, expectedPendingFebBalance = calculateMonthEndBalance(
				febTransactions,
				expectedAvailableFebBalance,
				expectedPendingFebBalance,
			)
			expectedPendingJanBalance := expectedPendingFebBalance
			expectedAvailableJanBalance := expectedAvailableFebBalance
			expectedAvailableJanBalance, expectedPendingJanBalance = calculateMonthEndBalance(
				janTransactions,
				expectedAvailableJanBalance,
				expectedPendingJanBalance,
			)
			expectedPendingDecemberBalance := expectedPendingJanBalance
			expectedAvailableDecemberBalance := expectedAvailableJanBalance
			expectedAvailableDecemberBalance, expectedPendingDecemberBalance = calculateMonthEndBalance(
				decemberTransactions,
				expectedAvailableDecemberBalance,
				expectedPendingDecemberBalance,
			)
			assert.Equal(t, 4, len(res.Months))
			assert.Equal(t, expectedAvailableMarchBalance.String(), monthVals[0])
			assert.Equal(t, expectedAvailableFebBalance.String(), monthVals[1])
			assert.Equal(t, expectedAvailableJanBalance.String(), monthVals[2])
			assert.Equal(t, expectedAvailableDecemberBalance.String(), monthVals[3])
		},
	)
}

func calculateMonthEndBalance(
	transactions []calculateTransferElement,
	expectedAvailableBalance decimal.Decimal,
	expectedPendingBalance decimal.Decimal,
) (decimal.Decimal, decimal.Decimal) {
	for _, transaction := range transactions {
		log.Printf("Transaction: %v", transaction.amount.String())
		if transaction.transactionType == model.Realized {
			expectedAvailableBalance = expectedAvailableBalance.Add(transaction.amount)
			expectedPendingBalance = expectedPendingBalance.Add(transaction.amount)
		} else {
			if transaction.status == model.Active {
				expectedPendingBalance = expectedPendingBalance.Add(transaction.amount)
			}
		}
	}
	return expectedAvailableBalance, expectedPendingBalance
}

func setupGetAccountHistoryTestCase(
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

func createRandomSimpleTransactions(
	firstAccountId string,
	secondAccountId string,
	startTime time.Time,
	endTime time.Time,
	numTransactions int,
) []interface{} {
	firstAccountObjectId, _ := pkgutils.StringToObjectId(firstAccountId)
	secondAccountObjectId, _ := pkgutils.StringToObjectId(secondAccountId)
	transactions := make([]interface{}, 0)

	for i := 0; i < numTransactions; i++ {
		randomAccount := rand.Intn(2)
		randomActiveOrPending := rand.Intn(2)
		var fromAccountId, toAccountId primitive.ObjectID

		if randomAccount == 0 {
			fromAccountId = firstAccountObjectId
			toAccountId = secondAccountObjectId
		} else {
			fromAccountId = secondAccountObjectId
			toAccountId = firstAccountObjectId
		}

		var activeOrPending, status string
		var expiryDate time.Time
		if randomActiveOrPending == 0 {
			activeOrPending = "realized"
		} else {
			activeOrPending = "pending"
			status = "active"
			expiryDate = time.Date(2075, time.March, 30, 0, 0, 0, 0, time.UTC)
		}

		var randomAmount = rand.Intn(11)
		var randomDecimals = strconv.Itoa(rand.Intn(100))

		amountDecimal, _ := primitive.ParseDecimal128(strconv.Itoa(randomAmount) + "." + randomDecimals)

		transaction := mongodb.MongoTransactionInput{
			FromBankAccountId: fromAccountId,
			ToBankAccountId:   toAccountId,
			Amount:            amountDecimal,
			Type:              activeOrPending,
			Status:            status,
			ExpirationDate:    pkgutils.TimeToTimestamp(expiryDate),
			CreatedAt:         pkgutils.TimeToTimestamp(randomTimeBetween(startTime, endTime)),
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func randomTimeBetween(start, end time.Time) time.Time {
	diff := end.Sub(start)
	randomDuration := time.Duration(rand.Int63n(int64(diff)))
	return start.Add(randomDuration)
}

type calculateTransferElement struct {
	amount          decimal.Decimal
	transactionType model.TransactionType
	status          model.PendingTransactionStatus
}
