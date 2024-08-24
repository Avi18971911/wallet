package dto

import "time"

// KnownAccountDTO represents an account known to and recognized by a particular account
// @swagger:model KnownAccountDTO
type KnownAccountDTO struct {
	// The account number of the known account
	// Required: true
	AccountNumber string `json:"accountNumber"`
	// The name of the account holder
	// Required: true
	AccountHolder string `json:"accountHolder"`
	// The type of the account (e.g., savings, checking)
	// Required: true
	AccountType string `json:"accountType"`
}

// AccountDetailsDTO represents the confidential details of an account belonging to a customer
// @swagger:model AccountDetailsDTO
type AccountDetailsDTO struct {
	// The unique identifier of the account
	// Required: true
	Id string `json:"id"`
	// The username associated with the account
	// Required: true
	Username string `json:"username"`
	// The available balance in the account
	// Required: true
	AvailableBalance float64 `json:"availableBalance"`
	// The account number
	// Required: true
	AccountNumber string `json:"accountNumber"`
	// The type of the account
	// Required: true
	AccountType string `json:"accountType"`
	// The account holder associated with the account
	// Required: true
	Person PersonDTO `json:"person"`
	// The list of accounts known to and recognized by the account holder
	// Required: true
	KnownAccounts []KnownAccountDTO `json:"knownAccounts"`
	// The creation timestamp of the account
	// Required: true
	CreatedAt time.Time `json:"createdAt"`
}

// AccountTransactionDTO represents a transaction between the given account and another account
// @swagger:model AccountTransactionDTO
type AccountTransactionDTO struct {
	// The unique identifier of the transaction
	// Required: true
	Id string `json:"id"`
	// The primary account ID associated with the transaction
	// Required: true
	AccountId string `json:"accountId"`
	// The other account ID involved in the transaction
	// Required: true
	OtherAccountId string `json:"otherAccountId"`
	// The type of the transaction (debit or credit)
	// Required: true
	TransactionType string `json:"transactionType"`
	// The amount involved in the transaction
	// Required: true
	Amount float64 `json:"amount"`
	// The timestamp of when the transaction was created
	// Required: true
	CreatedAt time.Time `json:"createdAt"`
}

// AccountLoginDTO represents the login credentials for an account
// @swagger:model AccountLoginDTO
type AccountLoginDTO struct {
	// The username for the login
	// Required: true
	Username string `json:"username"`
	// The password for the login
	// Required: true
	Password string `json:"password"`
}

// PersonDTO represents an account holder
// @swagger:model PersonDTO
type PersonDTO struct {
	// The first name of the person
	// Required: true
	FirstName string `json:"firstName"`
	// The last name of the person
	// Required: true
	LastName string `json:"lastName"`
}
