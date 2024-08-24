package dto

// TransactionRequest represents a request to add a new transaction from an account to another account.
// @swagger:model TransactionRequest
type TransactionRequest struct {
	// The account number of the account to which the amount is to be transferred
	// Required: true
	ToAccount string `json:"toAccount"`
	// The account number of the account from which the amount is to be transferred
	// Required: true
	FromAccount string `json:"fromAccount"`
	// The amount to be transferred
	// Required: true
	Amount float64 `json:"amount"`
}
