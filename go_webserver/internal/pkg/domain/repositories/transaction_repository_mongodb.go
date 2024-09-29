package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
)

type TransactionRepositoryMongodb struct {
	col *mongo.Collection
}

func CreateNewTransactionRepositoryMongodb(col *mongo.Collection) *TransactionRepositoryMongodb {
	ar := TransactionRepositoryMongodb{col: col}
	return &ar
}

func (tr *TransactionRepositoryMongodb) AddTransaction(
	details *model.TransactionDetails,
	ctx context.Context,
) error {
	mongoDetails, err := fromDomainTransactionDetails(details)
	if err != nil {
		return fmt.Errorf("error when converting domain TransactionDetails to mongo TransactionDetails "+
			"from BankAccount %s to BankAccount %s: %w", details.FromAccount, details.ToAccount, err)
	}
	_, err = tr.col.InsertOne(ctx, mongoDetails)
	if err != nil {
		return fmt.Errorf("error when inserting transaction from BankAccount %s to BankAccount %s: %w",
			details.FromAccount, details.ToAccount, err)
	}
	log.Printf("Successfully inserted transaction from "+
		"BankAccount %s to BankAccount %s\n", details.FromAccount, details.ToAccount)
	return nil
}

func (tr *TransactionRepositoryMongodb) GetAccountTransactions(
	accountId string, ctx context.Context,
) ([]model.AccountTransaction, error) {
	var res []model.AccountTransaction
	objectAccountId, err := utils.StringToObjectId(accountId)
	if err != nil {
		return nil, fmt.Errorf("error when converting account ID to object ID for accountId %s: %w", accountId, err)
	}
	pipeline := mongo.Pipeline{
		// Match transactions involving the accountId in either fromAccount or toAccount
		{{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"fromAccount", objectAccountId}},
				bson.D{{"toAccount", objectAccountId}},
			}},
		}}},
		// Add a new field 'transactionType' to indicate debit or credit transaction
		{{"$addFields", bson.D{
			{"transactionType", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromAccount", objectAccountId}}},
				"credit",
				"debit",
			}}}}},
		}},
		{{"$project", bson.D{
			{"_id", 1},
			{"_createdAt", 1},
			{"amount", 1},
			{"transactionType", 1},
			{"accountId", objectAccountId},
			{"otherAccountId", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromAccount", objectAccountId}}},
				"$toAccount",
				"$fromAccount",
			}}}},
		}}}}

	cursor, err := tr.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error when aggregating transactions for BankAccount %s: %w", accountId, err)
	}

	var mongoResults []mongodb.MongoAccountTransaction

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error when closing mongo Cursor when getting BankAccount Transactions "+
				"for BankAccount %s", accountId)
		}
	}()

	if err = cursor.All(ctx, &mongoResults); err != nil {
		return nil, fmt.Errorf("error when iterating over mongo Cursor when getting BankAccount Transactions "+
			"for BankAccount %s: %w", accountId, err)
	}
	if res, err = fromMongoAccountTransaction(mongoResults); err != nil {
		return nil, fmt.Errorf("error when converting mongo BankAccount Transactions to domain BankAccount "+
			"Transactions for BankAccount %s: %w", accountId, err)
	}
	log.Printf("Successfully retrieved BankAccount Transactions for BankAccount %s\n", accountId)
	return res, nil
}
