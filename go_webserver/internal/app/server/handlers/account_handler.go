package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/services"
)

type AccountDetailsDTO struct {
	Id               string    `json:"id"`
	Username         string    `json:"username"`
	AvailableBalance float64   `json:"availableBalance"`
	CreatedAt        time.Time `json:"createdAt"`
}

type AccountTransactionDTO struct {
	Id              string    `json:"id"`
	AccountId       string    `json:"accountId"`
	OtherAccountId  string    `json:"otherAccountId"`
	TransactionType string    `json:"transactionType"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"createdAt"`
}

func accountDetailsToDTO(tx *model.AccountDetails) AccountDetailsDTO {
	return AccountDetailsDTO{
		Id:               tx.Id,
		AvailableBalance: tx.AvailableBalance,
		Username:         tx.Username,
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
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{accountId} [get]
func AccountDetailsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountDetails, err := s.GetAccountDetails(accountID, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonAccountDetails := accountDetailsToDTO(accountDetails)
		err = json.NewEncoder(w).Encode(jsonAccountDetails)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
// @Failure 500 {string} string "Internal server error"
// @Router /accounts/{accountId}/transactions [get]
func AccountTransactionsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountTransactions, err := s.GetAccountTransactions(accountID, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonAccountTransactions := accountTransactionToDTO(accountTransactions)

		err = json.NewEncoder(w).Encode(jsonAccountTransactions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
