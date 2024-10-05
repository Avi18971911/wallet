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
	details *model.TransactionDetailsInput,
	ctx context.Context,
) error {
	mongoDetails, err := fromDomainTransactionDetails(details)
	if err != nil {
		return fmt.Errorf("error when converting domain TransactionDetailsInput to mongo TransactionDetailsInput "+
			"from BankAccount %s to BankAccount %s: %w", details.FromBankAccountId, details.ToBankAccountId, err)
	}
	_, err = tr.col.InsertOne(ctx, mongoDetails)
	if err != nil {
		return fmt.Errorf("error when inserting transaction from BankAccount %s to BankAccount %s: %w",
			details.FromBankAccountId, details.ToBankAccountId, err)
	}
	log.Printf("Successfully inserted transaction from "+
		"BankAccount %s to BankAccount %s\n", details.FromBankAccountId, details.ToBankAccountId)
	return nil
}

func (tr *TransactionRepositoryMongodb) GetTransactionsFromBankAccountId(
	bankAccountId string, ctx context.Context,
) ([]model.BankAccountTransactionOutput, error) {
	var res []model.BankAccountTransactionOutput
	objectAccountId, err := utils.StringToObjectId(bankAccountId)
	if err != nil {
		return nil, fmt.Errorf("error when converting account ID to object ID for "+
			"bankAccountId %s: %w", bankAccountId, err)
	}
	pipeline := mongo.Pipeline{
		// Match transactions involving the bankAccountId in either fromAccount or toAccount
		{{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"fromBankAccountId", objectAccountId}},
				bson.D{{"toBankAccountId", objectAccountId}},
			}},
		}}},
		// Add a new field 'transactionType' to indicate debit or credit transaction
		{{"$addFields", bson.D{
			{"transactionType", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromBankAccountId", objectAccountId}}},
				"credit",
				"debit",
			}}}}},
		}},
		{{"$project", bson.D{
			{"_id", 1},
			{"_createdAt", 1},
			{"amount", 1},
			{"transactionType", 1},
			{"bankAccountId", objectAccountId},
			{"otherBankAccountId", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromBankAccountId", objectAccountId}}},
				"$toBankAccountId",
				"$fromBankAccountId",
			}}}},
		}}}}

	cursor, err := tr.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error when aggregating transactions for BankAccount %s: %w", bankAccountId, err)
	}

	var mongoResults []mongodb.MongoAccountTransactionOutput

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error when closing mongo Cursor when getting BankAccount Transactions "+
				"for BankAccount %s", bankAccountId)
		}
	}()

	if err = cursor.All(ctx, &mongoResults); err != nil {
		return nil, fmt.Errorf("error when iterating over mongo Cursor when getting BankAccount Transactions "+
			"for BankAccount %s: %w", bankAccountId, err)
	}
	if res, err = fromMongoAccountTransaction(mongoResults); err != nil {
		return nil, fmt.Errorf("error when converting mongo BankAccount Transactions to domain BankAccount "+
			"Transactions for BankAccount %s: %w", bankAccountId, err)
	}
	log.Printf("Successfully retrieved BankAccount Transactions for BankAccount %s\n", bankAccountId)
	return res, nil
}
