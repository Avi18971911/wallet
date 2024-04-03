package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"webserver/services"
)

func AccountDetailsHandler(s services.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		s.GetAccountDetails(accountID)
		// TODO: Implement a return struct and encode into JSON
	}
}

func AccountTransactionsHandler(s services.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		s.GetAccountTransactions(accountID)
		// TODO: Implement a return struct and encode into JSON
	}
}
