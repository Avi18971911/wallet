package dto

// TransactionRequest represents a request to add a new transaction from an account to another account.
// @swagger:model TransactionRequest
type TransactionRequest struct {
	// The account number of the account to which the amount is to be transferred
	ToAccount string `json:"toAccount" validate:"required"`
	// The account number of the account from which the amount is to be transferred
	FromAccount string `json:"fromAccount" validate:"required"`
	// The amount to be transferred. Valid to two decimal places.
	Amount string `json:"amount" validate:"required"`
}
