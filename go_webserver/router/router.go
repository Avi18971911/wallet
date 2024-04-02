package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"webserver/handlers"
)

func createRouter(
	accountService handlers.AccountService,
	transactionService handlers.TransactionService,
) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/accounts/{accountID}", handlers.AccountDetailsHandler(accountService)).Methods("GET")
	r.Handle("/transactions", handlers.TransactionInsertHandler(transactionService)).Methods("POST")
	r.Handle("/accounts/{accountID}/transactions", handlers.AccountTransactionsHandler(accountService)).
		Methods("GET")

	return r
}
