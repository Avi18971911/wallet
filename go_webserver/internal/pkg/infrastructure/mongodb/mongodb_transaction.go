package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoTransactionInput struct {
	Id                primitive.ObjectID   `bson:"_id,omitempty"`
	FromBankAccountId primitive.ObjectID   `bson:"fromBankAccountId"`
	ToBankAccountId   primitive.ObjectID   `bson:"toBankAccountId"`
	Amount            primitive.Decimal128 `bson:"amount"`
	Type              string               `bson:"type"`
	ExpirationDate    primitive.Timestamp  `bson:"expirationDate,omitempty"`
	Status            string               `bson:"status,omitempty"`
	CreatedAt         primitive.Timestamp  `bson:"_createdAt"`
}
