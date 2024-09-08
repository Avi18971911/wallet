package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAccountDetails struct {
	Id            primitive.ObjectID  `bson:"_id"`
	Username      string              `bson:"username"`
	Password      string              `bson:"password"`
	Person        Person              `bson:"person"`
	Accounts      []Account           `bson:"accounts"`
	KnownAccounts []KnownAccount      `bson:"knownAccounts"`
	CreatedAt     primitive.Timestamp `bson:"_createdAt"`
}

type MongoAccountInput struct {
	Username      string              `bson:"username"`
	Password      string              `bson:"password"`
	Person        Person              `bson:"person"`
	Accounts      []Account           `bson:"accounts"`
	KnownAccounts []KnownAccount      `bson:"knownAccounts"`
	CreatedAt     primitive.Timestamp `bson:"_createdAt"`
}

type Person struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
}

type Account struct {
	Id               primitive.ObjectID `bson:"_id"`
	AccountNumber    string             `bson:"accountNumber"`
	AccountType      string             `bson:"accountType"`
	AvailableBalance float64            `bson:"availableBalance"`
}

type KnownAccount struct {
	Id            primitive.ObjectID `bson:"_id"`
	AccountNumber string             `bson:"accountNumber"`
	AccountHolder string             `bson:"accountHolder"`
	AccountType   string             `bson:"accountType"`
}

type MongoAccountTransaction struct {
	Id              primitive.ObjectID  `bson:"_id"`
	AccountId       primitive.ObjectID  `bson:"accountId"`
	OtherAccountId  primitive.ObjectID  `bson:"otherAccountId"`
	TransactionType string              `bson:"transactionType"`
	Amount          float64             `bson:"amount"`
	CreatedAt       primitive.Timestamp `bson:"_createdAt"`
}
