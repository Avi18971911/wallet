package router

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"webserver/handlers"
	"webserver/services"
)

func CreateRouter(
	accountService services.AccountService,
	transactionService services.TransactionService,
	ctx context.Context,
) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/accounts/{accountID}", handlers.AccountDetailsHandler(accountService, ctx)).Methods("GET")
	r.Handle("/transactions", handlers.TransactionInsertHandler(transactionService, ctx)).Methods("POST")
	r.Handle("/accounts/{accountID}/transactions", handlers.AccountTransactionsHandler(accountService, ctx)).
		Methods("GET")

	return r
}
