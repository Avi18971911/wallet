package dto

import (
	"time"
)

// KnownAccountDTO represents an account known to and recognized by a particular account
// @swagger:model KnownAccountDTO
type KnownAccountDTO struct {
	// The account ID of the known account
	Id string `json:"id" validate:"required"`
	// The account number of the known account
	AccountNumber string `json:"accountNumber" validate:"required"`
	// The name of the account holder
	AccountHolder string `json:"accountHolder" validate:"required"`
	// The type of the account (e.g., savings, checking)
	AccountType string `json:"accountType" validate:"required"`
}

// AccountDetailsDTO represents the confidential details of an account belonging to a customer
// @swagger:model AccountDetailsDTO
type AccountDetailsDTO struct {
	// The unique identifier of the account
	Id string `json:"id" validate:"required"`
	// The username associated with the account
	Username string `json:"username" validate:"required"`
	// The account holder associated with the account
	Person PersonDTO `json:"person" validate:"required"`
	// The list of accounts associated with the account holder
	Accounts []AccountDTO `json:"accounts" validate:"required"`
	// The list of accounts known to and recognized by the account holder
	KnownAccounts []KnownAccountDTO `json:"knownAccounts" validate:"required"`
	// The creation timestamp of the account
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

// AccountTransactionDTO represents a transaction between the given account and another account
// @swagger:model AccountTransactionDTO
type AccountTransactionDTO struct {
	// The unique identifier of the transaction
	Id string `json:"id" validate:"required"`
	// The primary account ID associated with the transaction
	AccountId string `json:"accountId" validate:"required"`
	// The other account ID involved in the transaction
	OtherAccountId string `json:"otherAccountId" validate:"required"`
	// The type of the transaction (debit or credit)
	TransactionType string `json:"transactionType" validate:"required"`
	// The amount involved in the transaction. Valid to two decimal places.
	Amount string `json:"amount" validate:"required"`
	// The timestamp of when the transaction was created
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

// AccountLoginDTO represents the login credentials for an account
// @swagger:model AccountLoginDTO
type AccountLoginDTO struct {
	// The username for the login
	Username string `json:"username" validate:"required"`
	// The password for the login
	Password string `json:"password" validate:"required"`
}

// AccountDTO represents an account associated with an account holder
// @swagger:model AccountDTO
type AccountDTO struct {
	// The unique identifier of the account
	Id string `json:"id" validate:"required"`
	// The account number associated with the account
	AccountNumber string `json:"accountNumber" validate:"required"`
	// The type of the account (e.g., savings, checking)
	AccountType string `json:"accountType" validate:"required"`
	// The available balance of the account. Valid to two decimal places.
	AvailableBalance string `json:"availableBalance" validate:"required"`
}

// PersonDTO represents an account holder
// @swagger:model PersonDTO
type PersonDTO struct {
	// The first name of the person
	FirstName string `json:"firstName" validate:"required"`
	// The last name of the person
	LastName string `json:"lastName" validate:"required"`
}
