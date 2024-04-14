package domain

type TransactionDetails struct {
	FromAccount string  `bson:"fromAccount"`
	ToAccount   string  `bson:"toAccount"`
	Amount      float64 `bson:"amount"`
}
