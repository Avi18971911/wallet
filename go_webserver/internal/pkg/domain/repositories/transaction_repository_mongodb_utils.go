package repositories

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
)

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
	decimal128Amount, err := utils.FromDecimalToPrimitiveDecimal128(details.Amount)
	if err != nil {
		return nil, fmt.Errorf("error when converting amount %s to Decimal128: %w", details.Amount, err)
	}
	return &mongodb.MongoTransactionDetails{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      decimal128Amount,
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
		decimalAmount, err := utils.FromPrimitiveDecimal128ToDecimal(elem.Amount)
		if err != nil {
			return res, fmt.Errorf("error when converting amount to decimal: %w", err)
		}
		res[i] = model.AccountTransaction{
			Id:              transactionId,
			AccountId:       accountId,
			OtherAccountId:  otherAccountId,
			TransactionType: elem.TransactionType,
			Amount:          decimalAmount,
			CreatedAt:       utils.TimestampToTime(elem.CreatedAt),
		}
	}
	return res, nil
}
