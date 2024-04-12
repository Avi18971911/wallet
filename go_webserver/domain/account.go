package domain

type AccountDetails struct {
	Id               string  `bson:"_id"`
	AvailableBalance float64 `bson:"availableBalance"`
}

type AccountTransaction struct {
	Id        string  `bson:"_id"`
	AccountId string  `bson:"accountId"`
	Amount    float64 `bson:"amount"`
	CreatedAt int64   `bson:"_createdAt"`
}
