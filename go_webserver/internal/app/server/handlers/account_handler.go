package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/infrastructure/mongodb"
)

type AccountDetailsDTO struct {
	Id               string  `json:"id"`
	AvailableBalance float64 `json:"availableBalance"`
}

type AccountTransactionDTO struct {
	Id        string    `json:"id"`
	AccountId string    `json:"accountId"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

func accountDetailsToDTO(tx *mongodb.MongoAccountDetails) AccountDetailsDTO {
	return AccountDetailsDTO{
		Id:               tx.Id,
		AvailableBalance: tx.AvailableBalance,
	}
}

func accountTransactionToDTO(tx []mongodb.MongoAccountTransaction) []AccountTransactionDTO {
	accountTransactionDTOList := make([]AccountTransactionDTO, len(tx))
	for i, element := range tx {
		accountTransactionDTOList[i] = AccountTransactionDTO{
			Id:        element.Id,
			AccountId: element.AccountId,
			Amount:    element.Amount,
			CreatedAt: element.CreatedAt,
		}
	}
	return accountTransactionDTOList
}

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
