package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"webserver/internal/app/server/dto"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/utils"
)

// TransactionInsertHandler creates a handler for adding a new transaction.
// @Summary Add a new transaction
// @Description Adds a new transaction to the system.
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body dto.TransactionRequest true "Transaction request"
// @Success 202 {string} string "Accepted"
// @Failure 400 {object} utils.ErrorMessage "Invalid request payload"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /transactions [post]
func TransactionInsertHandler(s services.TransactionService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.TransactionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.HttpError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Failed to close request body: %v", err)
			}
		}(r.Body)

		transactionInput, err := transactionDetailsToModel(&req)
		if err != nil {
			utils.HttpError(w, "Invalid amount given", http.StatusBadRequest)
			return
		}

		err = s.AddTransaction(transactionInput, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to Add Transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}
