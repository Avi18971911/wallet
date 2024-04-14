package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	domain2 "webserver/internal/pkg/domain"
)

type TransactionRepositoryMongodb struct {
	col *mongo.Collection
}

func CreateNewTransactionRepositoryMongodb(col *mongo.Collection) *TransactionRepositoryMongodb {
	ar := TransactionRepositoryMongodb{col: col}
	return &ar
}

func (tr *TransactionRepositoryMongodb) AddTransaction(
	details domain2.TransactionDetails,
	ctx context.Context,
) error {
	_, err := tr.col.InsertOne(ctx, details)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TransactionRepositoryMongodb) GetAccountTransactions(
	accountId string, ctx context.Context,
) ([]domain2.AccountTransaction, error) {
	pipeline := mongo.Pipeline{
		// Match transactions involving the accountId in either fromAccount or toAccount
		{{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"fromAccount", accountId}},
				bson.D{{"toAccount", accountId}},
			}},
		}}},
		// Add a new field 'type' to indicate debit or credit transaction
		{{"$addFields", bson.D{
			{"transactionType", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromAccount", accountId}}},
				"debit",
				"credit",
			}}}}},
		}},
		// Project desired fields, excluding the irrelevant account field
		{{"$project", bson.D{
			{"_id", 1},
			{"_createdAt", 1},
			{"type", 1},
			{"amount", 1},
			{"transactionType", 1},
			{"accountId", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$transactionType", "debit"}}},
				"$toAccount",
				"$fromAccount",
			}}}},
		}}}}

	cursor, err := tr.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []domain2.AccountTransaction

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error when closing mongo Cursor when getting Account Transactions "+
				"for Account %s", accountId)
		}
	}()

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
