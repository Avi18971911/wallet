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

func knownAccountToDTO(tx []model.KnownAccount) []dto.KnownAccountDTO {
	knownAccountDTOList := make([]dto.KnownAccountDTO, len(tx))
	for i, element := range tx {
		accountType, err := accountTypeEnumToString(element.AccountType)
		if err != nil {
			log.Printf("Failed to convert account type to string: %v", err)
		}
		knownAccountDTOList[i] = dto.KnownAccountDTO{
			AccountNumber: element.AccountNumber,
			AccountHolder: element.AccountHolder,
			AccountType:   accountType,
		}
	}
	return knownAccountDTOList
}

func accountsToDTO(tx []model.Account) []dto.AccountDTO {
	accountDTOList := make([]dto.AccountDTO, len(tx))
	for i, element := range tx {
		accountType, err := accountTypeEnumToString(element.AccountType)
		if err != nil {
			log.Printf("Failed to convert account type to string: %v", err)
		}
		accountDTOList[i] = dto.AccountDTO{
			AccountNumber:    element.AccountNumber,
			AccountType:      accountType,
			AvailableBalance: element.AvailableBalance,
		}
	}
	return accountDTOList
}

func accountDetailsToDTO(tx *model.AccountDetails) dto.AccountDetailsDTO {
	return dto.AccountDetailsDTO{
		Id:       tx.Id,
		Username: tx.Username,
		Person: dto.PersonDTO{
			FirstName: tx.Person.FirstName,
			LastName:  tx.Person.LastName,
		},
		Accounts:      accountsToDTO(tx.Accounts),
		KnownAccounts: knownAccountToDTO(tx.KnownAccounts),
		CreatedAt:     tx.CreatedAt,
	}
}

func accountTransactionToDTO(tx []model.AccountTransaction) []dto.AccountTransactionDTO {
	accountTransactionDTOList := make([]dto.AccountTransactionDTO, len(tx))
	for i, element := range tx {
		accountTransactionDTOList[i] = dto.AccountTransactionDTO{
			Id:              element.Id,
			AccountId:       element.AccountId,
			OtherAccountId:  element.OtherAccountId,
			TransactionType: element.TransactionType,
			Amount:          element.Amount,
			CreatedAt:       element.CreatedAt,
		}
	}
	return accountTransactionDTOList
}

func accountTypeEnumToString(tx int) (string, error) {
	switch tx {
	case model.Savings:
		return "savings", nil
	case model.Checking:
		return "checking", nil
	case model.Investment:
		return "investment", nil
	default:
		return "unknown", errors.New("invalid account type")
	}
}

// AccountDetailsHandler creates a handler for fetching account details.
// @Summary Get account details
// @Description Retrieves the details of a specific account by its ID.
// @Tags accounts
// @Accept json
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} dto.AccountDetailsDTO "Successful retrieval of account details"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/{accountId} [get]
func AccountDetailsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountDetails, err := s.GetAccountDetails(accountID, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to get Account Details", http.StatusInternalServerError)
			return
		}
		jsonAccountDetails := accountDetailsToDTO(accountDetails)
		err = json.NewEncoder(w).Encode(jsonAccountDetails)
		if err != nil {
			utils.HttpError(w, "Error encountered during response payload construction",
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// AccountTransactionsHandler creates a handler for fetching account transactions.
// @Summary Get account transactions
// @Description Retrieves a list of transactions for a specific account by its ID.
// @Tags transactions
// @Accept json
// @Produce json
// @Param accountId path string true "Account ID"
// @Success 200 {object} []dto.AccountTransactionDTO "Successful retrieval of account transactions"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/{accountId}/transactions [get]
func AccountTransactionsHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := mux.Vars(r)["accountId"]
		accountTransactions, err := s.GetAccountTransactions(accountID, ctx)
		if err != nil {
			utils.HttpError(w, "Failed to get Account Transactions", http.StatusInternalServerError)
			return
		}
		jsonAccountTransactions := accountTransactionToDTO(accountTransactions)

		err = json.NewEncoder(w).Encode(jsonAccountTransactions)
		if err != nil {
			utils.HttpError(w, "Error encountered during response payload construction",
				http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// AccountLoginHandler creates a handler for logging in a user.
// @Summary Login
// @Description Logs in a user with the provided username and password.
// @Tags accounts
// @Accept json
// @Produce json
// @Param login body dto.AccountLoginDTO true "Login payload"
// @Success 200 {object} dto.AccountDetailsDTO "Successful login"
// @Failure 401 {object} utils.ErrorMessage "Invalid credentials"
// @Failure 500 {object} utils.ErrorMessage "Internal server error"
// @Router /accounts/login [post]
func AccountLoginHandler(s services.AccountService, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.AccountLoginDTO
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

		w.WriteHeader(http.StatusOK)
	}
}
