package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/utils"
)

type AccountDetailsDTO struct {
	Id               string            `json:"id"`
	Username         string            `json:"username"`
	AvailableBalance float64           `json:"availableBalance"`
	AccountNumber    string            `json:"accountNumber"`
	Person           PersonDTO         `json:"person"`
	KnownAccounts    []KnownAccountDTO `json:"knownAccounts"`
	CreatedAt        time.Time         `json:"createdAt"`
}

type AccountTransactionDTO struct {
	Id              string    `json:"id"`
	AccountId       string    `json:"accountId"`
	OtherAccountId  string    `json:"otherAccountId"`
	TransactionType string    `json:"transactionType"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"createdAt"`
}

type AccountLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PersonDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type KnownAccountDTO struct {
	AccountNumber string `json:"accountNumber"`
	AccountHolder string `json:"accountHolder"`
	AccountType   int    `json:"accountType"`
}

func knownAccountToDTO(tx []model.KnownAccount) []KnownAccountDTO {
	knownAccountDTOList := make([]KnownAccountDTO, len(tx))
	for i, element := range tx {
		knownAccountDTOList[i] = KnownAccountDTO{
			AccountNumber: element.AccountNumber,
			AccountHolder: element.AccountHolder,
			AccountType:   element.AccountType,
		}
	}
	return knownAccountDTOList
}

func accountDetailsToDTO(tx *model.AccountDetails) AccountDetailsDTO {
	return AccountDetailsDTO{
		Id:               tx.Id,
		AvailableBalance: tx.AvailableBalance,
		Username:         tx.Username,
		AccountNumber:    tx.AccountNumber,
		Person: PersonDTO{
			FirstName: tx.Person.FirstName,
			LastName:  tx.Person.LastName,
		},
		KnownAccounts: knownAccountToDTO(tx.KnownAccounts),
		CreatedAt:     tx.CreatedAt,
	}
}

func accountTransactionToDTO(tx []model.AccountTransaction) []AccountTransactionDTO {
	accountTransactionDTOList := make([]AccountTransactionDTO, len(tx))
	for i, element := range tx {
		accountTransactionDTOList[i] = AccountTransactionDTO{
			Id:              element.Id,
			AccountId:       element.AccountId,
			OtherAccountId:  element.OtherAccountId,
			TransactionType: element.TransactionType,
			Amount:          element.Amount,
			CreatedAt:       element.CreatedAt,
		}
	}
	return accountTransactionDTOList
}

// AccountDetailsHandler creates a handler for fetching account details.
// @Summary Get account details
// @Description Retrieves the details of a specific account by its ID.
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} AccountDetailsDTO "Successful retrieval of account details"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/{accountId} [get]
func AccountDetailsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountDetails, err := s.GetAccountDetails(accountID, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to get Account Details", http.StatusInternalServerError)
			return
		}
		jsonAccountDetails := accountDetailsToDTO(accountDetails)
		err = json.NewEncoder(w).Encode(jsonAccountDetails)
		if err != nil {
			utils.HttpError(w, "Error encountered during response payload construction",
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// AccountTransactionsHandler creates a handler for fetching account transactions.
// @Summary Get account transactions
// @Description Retrieves a list of transactions for a specific account by its ID.
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} []AccountTransactionDTO "Successful retrieval of account transactions"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/{accountId}/transactions [get]
func AccountTransactionsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountTransactions, err := s.GetAccountTransactions(accountID, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to get Account Transactions", http.StatusInternalServerError)
			return
		}
		jsonAccountTransactions := accountTransactionToDTO(accountTransactions)

		err = json.NewEncoder(w).Encode(jsonAccountTransactions)
		if err != nil {
			utils.HttpError(w, "Error encountered during response payload construction",
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// AccountLoginHandler creates a handler for logging in a user.
// @Summary Login
// @Description Logs in a user with the provided username and password.
// @Tags accounts
// @Accept json
// @Produce json
// @Param login body AccountLoginDTO true "Login payload"
// @Success 200 {object} AccountDetailsDTO "Successful login"
// @Failure 401 {object} utils.ErrorMessage "Invalid credentials"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/login [post]
func AccountLoginHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AccountLoginDTO
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.HttpError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close request body: %v", err)
			}
		}(r.Body)

		accountDetails, err := s.Login(req.Username, req.Password, ctx)
		if err != nil {
			if errors.Is(err, model.ErrInvalidCredentials) {
				utils.HttpError(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			utils.HttpError(w, "Error encountered during login", http.StatusInternalServerError)
			return
		}

		jsonAccountDetails := accountDetailsToDTO(accountDetails)
		err = json.NewEncoder(w).Encode(jsonAccountDetails)
		if err != nil {
			utils.HttpError(w, "Error encountered during JSON Encoding of Response",
				http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
