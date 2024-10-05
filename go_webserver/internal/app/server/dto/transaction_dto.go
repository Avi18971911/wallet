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

// TransactionsForBankAccountRequest represents a request to retrieve transactions for a bank account.
// @swagger:model TransactionsForBankAccountRequest
type TransactionsForBankAccountRequest struct {
	// The bank account ID of the account for which transactions are to be retrieved
	BankAccountId string `json:"bankAccountId" validate:"required"`
	// The earliest time (inclusive) from which transactions are to be retrieved. Format: RFC3339
	FromTime string `json:"fromTime"`
	// The latest time (inclusive) from which transactions are to be retrieved. Format: RFC3339
	ToTime string `json:"toTime"`
}
