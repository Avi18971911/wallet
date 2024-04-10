package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"webserver/services"
)

func AccountDetailsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		s.GetAccountDetails(accountID, ctx)
		// TODO: Implement a return struct and encode into JSON
	}
}

func AccountTransactionsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		s.GetAccountTransactions(accountID, ctx)
		// TODO: Implement a return struct and encode into JSON
	}
}
