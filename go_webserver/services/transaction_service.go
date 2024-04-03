package services

type TransactionService interface {
	AddTransaction(toAccount string, fromAccount string, amount float64)
}
