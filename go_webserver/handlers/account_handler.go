package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type AccountService interface {
	getAccountDetails(accountId string)
	getAccountTransactions(accountId string)
}

func AccountDetailsHandler(s AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		s.getAccountDetails(accountID)
		// TODO: Implement a return struct and encode into JSON
	}
}

func AccountTransactionsHandler(s AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		s.getAccountTransactions(accountID)
		// TODO: Implement a return struct and encode into JSON
	}
}
