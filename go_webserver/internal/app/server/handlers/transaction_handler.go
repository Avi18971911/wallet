package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"webserver/internal/pkg/domain/services"
)

type TransactionRequest struct {
	ToAccount   string  `json:"toAccount"`
	FromAccount string  `json:"fromAccount"`
	Amount      float64 `json:"amount"`
}

// TransactionInsertHandler creates a handler for adding a new transaction.
// @Summary Add a new transaction
// @Description Adds a new transaction to the system.
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body TransactionRequest true "Transaction request"
// @Success 202 {string} string "Accepted"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /transactions [post]
func TransactionInsertHandler(s services.TransactionService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransactionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close request body: %v", err)
			}
		}(r.Body)

		err = s.AddTransaction(req.ToAccount, req.FromAccount, req.Amount, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
