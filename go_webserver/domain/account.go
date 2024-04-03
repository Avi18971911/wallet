package domain

type AccountDetails struct {
	Id               string
	AvailableBalance float64
}

type AccountTransaction struct {
	Id        string
	AccountId string
	Amount    float64
	CreatedAt int64
}
