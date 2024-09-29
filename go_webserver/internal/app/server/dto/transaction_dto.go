package dto

// TransactionRequest represents a request to add a new transaction from an account to another account.
// @swagger:model TransactionRequest
type TransactionRequest struct {
	// The bank account ID of the account to which the amount is to be transferred
	ToBankAccountId string `json:"toBankAccountId" validate:"required"`
	// The bank account ID of the account from which the amount is to be transferred
	FromBankAccountId string `json:"fromBankAccountId" validate:"required"`
	// The amount to be transferred. Valid to two decimal places.
	Amount string `json:"amount" validate:"required"`
}
