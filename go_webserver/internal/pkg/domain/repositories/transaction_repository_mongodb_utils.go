package repositories

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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

	expirationDate := utils.TimeToTimestamp(details.ExpirationDate)

	transactionType, err := getStringFromTransactionType(details.Type)
	if err != nil {
		log.Printf("Error when converting transaction type: %v", err)
		return nil, err
	}

	status := getStringFromPendingTransactionStatus(details.Status)
	if err != nil {
		log.Printf("Error when converting pending transaction status: %v", err)
		return nil, err
	}

	return &mongodb.MongoTransactionInput{
		FromBankAccountId: fromAccount,
		ToBankAccountId:   toAccount,
		Amount:            decimal128Amount,
		CreatedAt:         utils.GetCurrentTimestamp(),
		Type:              transactionType,
		ExpirationDate:    expirationDate,
		Status:            status,
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

func getStringFromTransactionType(transactionType model.TransactionType) (string, error) {
	switch transactionType {
	case model.Realized:
		return "realized", nil
	case model.Pending:
		return "pending", nil
	default:
		return "unknown", fmt.Errorf("unknown transaction type: %s", transactionType)
	}
}

func getStringFromPendingTransactionStatus(status model.PendingTransactionStatus) string {
	switch status {
	case model.Active:
		return "active"
	case model.Applied:
		return "applied"
	case model.Revoked:
		return "revoked"
	default:
		return ""
	}
}
