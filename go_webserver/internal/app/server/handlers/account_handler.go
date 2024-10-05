package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"webserver/internal/app/server/dto"
	"webserver/internal/pkg/domain/model"
	"webserver/internal/pkg/domain/services"
	"webserver/internal/pkg/utils"
)

// AccountDetailsHandler creates a handler for fetching account details.
// @Summary Get account details
// @Description Retrieves the details of a specific account by its ID.
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountId path string true "BankAccount ID"
// @Success 200 {object} dto.AccountDetailsResponseDTO "Successful retrieval of account details"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/{accountId} [get]
func AccountDetailsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountDetails, err := s.GetAccountDetailsFromBankAccountId(accountID, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to get BankAccount Details", http.StatusInternalServerError)
			return
		}
		jsonAccountDetails := accountDetailsToDTO(accountDetails)
		err = json.NewEncoder(w).Encode(jsonAccountDetails)
		if err != nil {
			utils.HttpError(w, "Error encountered during response payload construction",
				http.StatusInternalServerError)
			return
		}
	}
}

// AccountTransactionsHandler creates a handler for fetching account transactions.
// @Summary Get account transactions
// @Description Retrieves a list of transactions for a specific account by its ID.
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountId path string true "BankAccount ID"
// @Success 200 {object} []dto.AccountTransactionResponseDTO "Successful retrieval of account transactions"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/{accountId}/transactions [get]
func AccountTransactionsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.AccountTransactionRequestDTO
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.HttpError(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		accountTransactionsInput := accountTransactionRequestToInput(req)
		accountTransactions, err := s.GetBankAccountTransactions(&accountTransactionsInput, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to get BankAccount Transactions", http.StatusInternalServerError)
			return
		}
		jsonAccountTransactions := accountTransactionToDTO(accountTransactions)

		err = json.NewEncoder(w).Encode(jsonAccountTransactions)
		if err != nil {
			utils.HttpError(w, "Error encountered during response payload construction",
				http.StatusInternalServerError)
			return
		}
	}
}

// AccountLoginHandler creates a handler for logging in a user.
// @Summary Login
// @Description Logs in a user with the provided username and password.
// @Tags accounts
// @Accept json
// @Produce json
// @Param login body dto.AccountLoginRequestDTO true "Login payload"
// @Success 200 {object} dto.AccountDetailsResponseDTO "Successful login"
// @Failure 401 {object} utils.ErrorMessage "Invalid credentials"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/login [post]
func AccountLoginHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.AccountLoginRequestDTO
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

		accountDetails, err := s.Login(req.Username, req.Password, ctx)
		if err != nil {
			if errors.Is(err, model.ErrInvalidCredentials) {
				utils.HttpError(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			utils.HttpError(w, "Error encountered during login", http.StatusInternalServerError)
			return
		}

		jsonAccountDetails := accountDetailsToDTO(accountDetails)
		err = json.NewEncoder(w).Encode(jsonAccountDetails)
		if err != nil {
			utils.HttpError(w, "Error encountered during JSON Encoding of Response",
				http.StatusInternalServerError)
			return
		}

	}
}
