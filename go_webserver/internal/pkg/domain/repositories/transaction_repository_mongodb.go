package repositories

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			"from Account %s to Account %s: %w", details.FromAccount, details.ToAccount, err)
	}
	_, err = tr.col.InsertOne(ctx, mongoDetails)
	if err != nil {
		return fmt.Errorf("error when inserting transaction from Account %s to Account %s: %w",
			details.FromAccount, details.ToAccount, err)
	}
	log.Printf("Successfully inserted transaction from "+
		"Account %s to Account %s\n", details.FromAccount, details.ToAccount)
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
		return nil, fmt.Errorf("error when aggregating transactions for Account %s: %w", accountId, err)
	}

	var mongoResults []mongodb.MongoAccountTransaction

	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("Error when closing mongo Cursor when getting Account Transactions "+
				"for Account %s", accountId)
		}
	}()

	if err = cursor.All(ctx, &mongoResults); err != nil {
		return nil, fmt.Errorf("error when iterating over mongo Cursor when getting Account Transactions "+
			"for Account %s: %w", accountId, err)
	}
	if res, err = fromMongoAccountTransaction(mongoResults); err != nil {
		return nil, fmt.Errorf("error when converting mongo Account Transactions to domain Account "+
			"Transactions for Account %s: %w", accountId, err)
	}
	log.Printf("Successfully retrieved Account Transactions for Account %s\n", accountId)
	return res, nil
}

func fromDomainTransactionDetails(details *model.TransactionDetails) (*mongodb.MongoTransactionDetails, error) {
	var fromAccount, toAccount primitive.ObjectID
	var err error
	fromAccount, err = utils.StringToObjectId(details.FromAccount)
	if err != nil {
		return nil, fmt.Errorf("error when converting fromAccount %s to ObjectID: %w", details.FromAccount, err)
	}
	toAccount, err = utils.StringToObjectId(details.ToAccount)
	if err != nil {
		return nil, fmt.Errorf("error when converting toAccount %s to ObjectID: %w", details.ToAccount, err)
	}
	return &mongodb.MongoTransactionDetails{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      details.Amount,
		CreatedAt:   utils.GetCurrentTimestamp(),
	}, nil
}

func fromMongoAccountTransaction(
	accountTransactions []mongodb.MongoAccountTransaction,
) ([]model.AccountTransaction, error) {
	var res = make([]model.AccountTransaction, len(accountTransactions))
	var err error
	var transactionId, accountId, otherAccountId string
	for i, elem := range accountTransactions {
		transactionId, err = utils.ObjectIdToString(elem.Id)
		if err != nil {
			return res, fmt.Errorf("error when converting transaction ID to string: %w", err)
		}
		accountId, err = utils.ObjectIdToString(elem.AccountId)
		if err != nil {
			return res, fmt.Errorf("error when converting account ID to string: %w", err)
		}
		otherAccountId, err = utils.ObjectIdToString(elem.OtherAccountId)
		if err != nil {
			return res, fmt.Errorf("error when converting other account ID to string: %w", err)
		}
		res[i] = model.AccountTransaction{
			Id:              transactionId,
			AccountId:       accountId,
			OtherAccountId:  otherAccountId,
			TransactionType: elem.TransactionType,
			Amount:          elem.Amount,
			CreatedAt:       utils.TimestampToTime(elem.CreatedAt),
		}
	}
	return res, nil
}
