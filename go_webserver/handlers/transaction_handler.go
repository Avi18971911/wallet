package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type TransactionRequest struct {
	ToAccount   string  `json:"toAccount"`
	FromAccount string  `json:"fromAccount"`
	Amount      float64 `json:"amount"`
}

type TransactionService interface {
	addTransaction(toAccount string, fromAccount string, amount float64)
}

func TransactionInsertHandler(s TransactionService) http.HandlerFunc {
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

		s.addTransaction(req.ToAccount, req.FromAccount, req.Amount)
		w.WriteHeader(http.StatusAccepted)
		// TODO: Think if anything else is required
	}
}
