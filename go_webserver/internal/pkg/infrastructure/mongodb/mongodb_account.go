package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAccountDetails struct {
	Id               primitive.ObjectID  `bson:"_id"`
	AvailableBalance float64             `bson:"availableBalance"`
	Username         string              `bson:"username"`
	Password         string              `bson:"password"`
	Person           Person              `bson:"person"`
	AccountNumber    string              `bson:"accountNumber"`
	AccountType      string              `bson:"accountType"`
	KnownAccounts    []KnownAccount      `bson:"knownAccounts"`
	CreatedAt        primitive.Timestamp `bson:"_createdAt"`
}

type MongoAccountInput struct {
	Username        string              `bson:"username"`
	Password        string              `bson:"password"`
	AccountNumber   string              `bson:"accountNumber"`
	AccountType     string              `bson:"accountType"`
	StartingBalance float64             `bson:"availableBalance"`
	Person          Person              `bson:"person"`
	KnownAccounts   []KnownAccount      `bson:"knownAccounts"`
	CreatedAt       primitive.Timestamp `bson:"_createdAt"`
}

type Person struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}

type KnownAccount struct {
	AccountNumber string `bson:"accountNumber"`
	AccountHolder string `bson:"accountHolder"`
	AccountType   string `bson:"accountType"`
}

type MongoAccountTransaction struct {
	Id              primitive.ObjectID  `bson:"_id"`
	AccountId       primitive.ObjectID  `bson:"accountId"`
	OtherAccountId  primitive.ObjectID  `bson:"otherAccountId"`
	TransactionType string              `bson:"transactionType"`
	Amount          float64             `bson:"amount"`
	CreatedAt       primitive.Timestamp `bson:"_createdAt"`
}
