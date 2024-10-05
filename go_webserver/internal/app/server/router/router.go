package router

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"webserver/internal/app/server/handlers"
	"webserver/internal/pkg/domain/services"
)

func CreateRouter(
	accountService services.AccountService,
	transactionService services.TransactionService,
	ctx context.Context,
) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/accounts/{accountId}", handlers.AccountDetailsHandler(accountService, ctx)).Methods("GET")
	r.Handle("/transactions", handlers.TransactionInsertHandler(transactionService, ctx)).Methods("POST")
	r.Handle("/accounts/transactions", handlers.AccountTransactionsHandler(accountService, ctx)).Methods("GET")
	r.Handle("/accounts/login", handlers.AccountLoginHandler(accountService, ctx)).Methods("POST")
	return r
}
