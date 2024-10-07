package repositories

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
	transactionType := string(details.Type)
	status := string(details.Status)

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
	accountTransactions []mongodb.MongoAccountTransactionOutput,
) ([]model.BankAccountTransactionOutput, error) {
	var res = make([]model.BankAccountTransactionOutput, len(accountTransactions))
	var err error
	var transactionId, accountId, otherAccountId string
	for i, elem := range accountTransactions {
		transactionId, err = utils.ObjectIdToString(elem.Id)
		if err != nil {
			return res, fmt.Errorf(
				"error when converting transaction ID to string for accountId %s: %w", accountId, err,
			)
		}
		accountId, err = utils.ObjectIdToString(elem.BankAccountId)
		if err != nil {
			return res, fmt.Errorf(
				"error when converting account ID to string for accountId %s: %w", accountId, err,
			)
		}
		otherAccountId, err = utils.ObjectIdToString(elem.OtherBankAccountId)
		if err != nil {
			return res, fmt.Errorf(
				"error when converting other account ID to string for accountId %s: %w", accountId, err,
			)
		}
		decimalAmount, err := utils.FromPrimitiveDecimal128ToDecimal(elem.Amount)
		if err != nil {
			return res, fmt.Errorf(
				"error when converting amount to decimal for accountId %s: %w", accountId, err,
			)
		}
		transactionType, err := toTransactionType(elem.Type)
		if err != nil {
			return res, fmt.Errorf("error when parsing transactionType for accountId %s: %w", accountId, err)
		}
		var pendingTransactionStatus model.PendingTransactionStatus
		if transactionType == model.Pending {
			pendingTransactionStatus, err = toPendingTransactionStatus(elem.Status)
		} else {
			pendingTransactionStatus = ""
		}
		if err != nil {
			return res, fmt.Errorf("error when parsing pendingTransactionStatus for accountId %s: %w", accountId, err)
		}
		var expirationDateStamp time.Time
		if transactionType == model.Pending {
			expirationDateStamp = utils.TimestampToTime(elem.ExpirationDate)
		}
		res[i] = model.BankAccountTransactionOutput{
			Id:                 transactionId,
			BankAccountId:      accountId,
			OtherBankAccountId: otherAccountId,
			TransactionType:    transactionType,
			ExpirationDate:     expirationDateStamp,
			Status:             pendingTransactionStatus,
			Amount:             decimalAmount,
			CreatedAt:          utils.TimestampToTime(elem.CreatedAt),
		}
	}
	return res, nil
}

func fromDomainTransactionForBankAccountInput(
	input *model.TransactionsForBankAccountInput,
) (*mongodb.MongoTransactionForBankAccountInput, error) {
	bankAccountId, err := utils.StringToObjectId(input.BankAccountId)
	if err != nil {
		return nil, fmt.Errorf("error when converting bank account ID to ObjectID: %w", err)
	}
	return &mongodb.MongoTransactionForBankAccountInput{
		BankAccountId: bankAccountId,
		FromTime:      utils.TimeToTimestamp(input.FromTime),
		ToTime:        utils.TimeToTimestamp(input.ToTime),
	}, nil
}

func toTransactionType(transactionType string) (model.TransactionType, error) {
	switch transactionType {
	case "realized":
		return model.Realized, nil
	case "pending":
		return model.Pending, nil
	default:
		return "", fmt.Errorf("invalid transaction type: %s", transactionType)
	}
}

func toPendingTransactionStatus(status string) (model.PendingTransactionStatus, error) {
	switch status {
	case "active":
		return model.Active, nil
	case "applied":
		return model.Applied, nil
	case "revoked":
		return model.Revoked, nil
	default:
		return "", fmt.Errorf("invalid pending transaction status: %s", status)
	}
}
