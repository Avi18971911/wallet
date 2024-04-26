package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoTransactionDetails struct {
	FromAccount primitive.ObjectID  `bson:"fromAccount"`
	ToAccount   primitive.ObjectID  `bson:"toAccount"`
	Amount      float64             `bson:"amount"`
	CreatedAt   primitive.Timestamp `bson:"_createdAt"`
}
