package router

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"webserver/internal/app/server/handlers"
	services2 "webserver/internal/pkg/services"
)

func CreateRouter(
	accountService services2.AccountService,
	transactionService services2.TransactionService,
	ctx context.Context,
) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/accounts/{accountID}", handlers.AccountDetailsHandler(accountService, ctx)).Methods("GET")
	r.Handle("/transactions", handlers.TransactionInsertHandler(transactionService, ctx)).Methods("POST")
	r.Handle("/accounts/{accountID}/transactions", handlers.AccountTransactionsHandler(accountService, ctx)).
		Methods("GET")

	return r
}
