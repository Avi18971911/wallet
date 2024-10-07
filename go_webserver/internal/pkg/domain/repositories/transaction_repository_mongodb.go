package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
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
	input *model.TransactionsForBankAccountInput,
	ctx context.Context,
) ([]model.BankAccountTransactionOutput, error) {
	var res []model.BankAccountTransactionOutput
	mongoInput, err := fromDomainTransactionForBankAccountInput(input)
	if err != nil {
		return nil, fmt.Errorf("error when converting domain TransactionsForBankAccountInput to mongo "+
			"TransactionsForBankAccountInput for BankAccount %s: %w", input.BankAccountId, err)
	}
	pipeline := mongo.Pipeline{
		// Match transactions involving the bankAccountId in either fromAccount or toAccount
		{{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"fromBankAccountId", mongoInput.BankAccountId}},
				bson.D{{"toBankAccountId", mongoInput.BankAccountId}},
			}},
			{"_createdAt", bson.D{
				{"$gte", mongoInput.FromTime},
				{"$lte", mongoInput.ToTime},
			}},
		}}},
		// Add a new field 'transactionType' to indicate debit or credit transaction
		{{"$addFields", bson.D{
			{"transactionNature", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromBankAccountId", mongoInput.BankAccountId}}},
				"credit",
				"debit",
			}}}}},
		}},
		{{"$project", bson.D{
			{"_id", 1},
			{"_createdAt", 1},
			{"amount", 1},
			{"transactionNature", 1},
			{"type", 1},
			{"status", 1},
			{"expirationDate", 1},
			{"bankAccountId", mongoInput.BankAccountId},
			{"otherBankAccountId", bson.D{{"$cond", bson.A{
				bson.D{{"$eq", bson.A{"$fromBankAccountId", mongoInput.BankAccountId}}},
				"$toBankAccountId",
				"$fromBankAccountId",
			}}}},
		}}}}

	cursor, err := tr.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("error when aggregating transactions for BankAccount %s: %w", input.BankAccountId, err)
	}

	var mongoResults []mongodb.MongoAccountTransactionOutput

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error when closing mongo Cursor when getting BankAccount Transactions "+
				"for BankAccount %s", input.BankAccountId)
		}
	}()

	if err = cursor.All(ctx, &mongoResults); err != nil {
		return nil, fmt.Errorf("error when iterating over mongo Cursor when getting BankAccount Transactions "+
			"for BankAccount %s: %w", input.BankAccountId, err)
	}
	if res, err = fromMongoAccountTransaction(mongoResults); err != nil {
		return nil, fmt.Errorf("error when converting mongo BankAccount Transactions to domain BankAccount "+
			"Transactions for BankAccount %s: %w", input.BankAccountId, err)
	}
	log.Printf("Successfully retrieved BankAccount Transactions for BankAccount %s\n", input.BankAccountId)
	return res, nil
}
