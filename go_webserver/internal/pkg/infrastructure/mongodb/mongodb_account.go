package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MongoAccountDetails struct {
	Id               primitive.ObjectID `bson:"_id"`
	AvailableBalance float64            `bson:"availableBalance"`
}

type MongoAccountTransaction struct {
	Id              primitive.ObjectID `bson:"_id"`
	AccountId       primitive.ObjectID `bson:"accountId"`
	OtherAccountId  primitive.ObjectID `bson:"otherAccountId"`
	TransactionType string             `bson:"transactionType"`
	Amount          float64            `bson:"amount"`
	CreatedAt       time.Time          `bson:"_createdAt"`
}
