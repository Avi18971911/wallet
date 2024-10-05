package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoTransactionInput struct {
	FromBankAccountId primitive.ObjectID   `bson:"fromBankAccountId"`
	ToBankAccountId   primitive.ObjectID   `bson:"toBankAccountId"`
	Amount            primitive.Decimal128 `bson:"amount"`
	Type              string               `bson:"type"`
	ExpirationDate    primitive.Timestamp  `bson:"expirationDate,omitempty"`
	Status            string               `bson:"status,omitempty"`
	CreatedAt         primitive.Timestamp  `bson:"_createdAt"`
}

type MongoTransactionForBankAccountInput struct {
	BankAccountId primitive.ObjectID  `bson:"bankAccountId"`
	FromTime      primitive.Timestamp `bson:"fromTime"`
	ToTime        primitive.Timestamp `bson:"toTime"`
}
