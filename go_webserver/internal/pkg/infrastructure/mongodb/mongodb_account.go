package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoAccountDetails struct {
	Id               primitive.ObjectID  `bson:"_id"`
	AvailableBalance float64             `bson:"availableBalance"`
	Username         string              `bson:"username"`
	Password         string              `bson:"password"`
	CreatedAt        primitive.Timestamp `bson:"_createdAt"`
}

type MongoAccountTransaction struct {
	Id              primitive.ObjectID  `bson:"_id"`
	AccountId       primitive.ObjectID  `bson:"accountId"`
	OtherAccountId  primitive.ObjectID  `bson:"otherAccountId"`
	TransactionType string              `bson:"transactionType"`
	Amount          float64             `bson:"amount"`
	CreatedAt       primitive.Timestamp `bson:"_createdAt"`
}
