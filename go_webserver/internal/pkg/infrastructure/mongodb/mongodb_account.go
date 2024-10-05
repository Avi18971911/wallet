package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAccountOutput struct {
	Id                primitive.ObjectID  `bson:"_id"`
	Username          string              `bson:"username"`
	Password          string              `bson:"password"`
	Person            Person              `bson:"person"`
	BankAccounts      []BankAccount       `bson:"bankAccounts"`
	KnownBankAccounts []KnownBankAccount  `bson:"knownBankAccounts"`
	CreatedAt         primitive.Timestamp `bson:"_createdAt"`
}

type MongoAccountInput struct {
	Username          string              `bson:"username"`
	Password          string              `bson:"password"`
	Person            Person              `bson:"person"`
	BankAccounts      []BankAccount       `bson:"bankAccounts"`
	KnownBankAccounts []KnownBankAccount  `bson:"knownBankAccounts"`
	CreatedAt         primitive.Timestamp `bson:"_createdAt"`
}

type Person struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}

type BankAccount struct {
	Id               primitive.ObjectID   `bson:"_id,omitempty"`
	AccountNumber    string               `bson:"accountNumber"`
	AccountType      string               `bson:"accountType"`
	PendingBalance   primitive.Decimal128 `bson:"pendingBalance"`
	AvailableBalance primitive.Decimal128 `bson:"availableBalance"`
}

type KnownBankAccount struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	AccountNumber string             `bson:"accountNumber"`
	AccountHolder string             `bson:"accountHolder"`
	AccountType   string             `bson:"accountType"`
}

type MongoAccountTransactionOutput struct {
	Id                 primitive.ObjectID   `bson:"_id"`
	BankAccountId      primitive.ObjectID   `bson:"bankAccountId"`
	OtherBankAccountId primitive.ObjectID   `bson:"otherBankAccountId"`
	TransactionType    string               `bson:"transactionType"`
	Amount             primitive.Decimal128 `bson:"amount"`
	CreatedAt          primitive.Timestamp  `bson:"_createdAt"`
}
