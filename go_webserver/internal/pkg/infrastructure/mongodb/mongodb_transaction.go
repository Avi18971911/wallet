package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoTransactionDetails struct {
	Id          primitive.ObjectID   `bson:"_id"`
	FromAccount primitive.ObjectID   `bson:"fromAccount"`
	ToAccount   primitive.ObjectID   `bson:"toAccount"`
	Amount      primitive.Decimal128 `bson:"amount"`
	CreatedAt   primitive.Timestamp  `bson:"_createdAt"`
}
