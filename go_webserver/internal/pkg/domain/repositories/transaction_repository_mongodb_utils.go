package repositories

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/infrastructure/mongodb"
	"webserver/internal/pkg/utils"
)

func fromDomainTransactionDetails(details *model.TransactionDetailsInput) (*mongodb.MongoTransactionInput, error) {
	var fromAccount, toAccount primitive.ObjectID
	var err error
	fromAccount, err = utils.StringToObjectId(details.FromBankAccountId)
	if err != nil {
		return nil, fmt.Errorf("error when converting fromAccount %s to ObjectID: %w", details.FromBankAccountId, err)
	}
	toAccount, err = utils.StringToObjectId(details.ToBankAccountId)
	if err != nil {
		return nil, fmt.Errorf("error when converting toAccount %s to ObjectID: %w", details.ToBankAccountId, err)
	}
	decimal128Amount, err := utils.FromDecimalToPrimitiveDecimal128(details.Amount)
	if err != nil {
		return nil, fmt.Errorf("error when converting amount %s to Decimal128: %w", details.Amount, err)
	}
	return &mongodb.MongoTransactionInput{
		FromBankAccountId: fromAccount,
		ToBankAccountId:   toAccount,
		Amount:            decimal128Amount,
		CreatedAt:         utils.GetCurrentTimestamp(),
	}, nil
}

func fromMongoAccountTransaction(
	accountTransactions []mongodb.MongoAccountTransaction,
) ([]model.BankAccountTransaction, error) {
	var res = make([]model.BankAccountTransaction, len(accountTransactions))
	var err error
	var transactionId, accountId, otherAccountId string
	for i, elem := range accountTransactions {
		transactionId, err = utils.ObjectIdToString(elem.Id)
		if err != nil {
			return res, fmt.Errorf("error when converting transaction ID to string: %w", err)
		}
		accountId, err = utils.ObjectIdToString(elem.BankAccountId)
		if err != nil {
			return res, fmt.Errorf("error when converting account ID to string: %w", err)
		}
		otherAccountId, err = utils.ObjectIdToString(elem.OtherBankAccountId)
		if err != nil {
			return res, fmt.Errorf("error when converting other account ID to string: %w", err)
		}
		decimalAmount, err := utils.FromPrimitiveDecimal128ToDecimal(elem.Amount)
		if err != nil {
			return res, fmt.Errorf("error when converting amount to decimal: %w", err)
		}
		res[i] = model.BankAccountTransaction{
			Id:                 transactionId,
			BankAccountId:      accountId,
			OtherBankAccountId: otherAccountId,
			TransactionType:    elem.TransactionType,
			Amount:             decimalAmount,
			CreatedAt:          utils.TimestampToTime(elem.CreatedAt),
		}
	}
	return res, nil
}
