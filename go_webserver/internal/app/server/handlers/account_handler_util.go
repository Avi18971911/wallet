package handlers

import (
	"errors"
	"log"
	"webserver/internal/app/server/dto"
	"webserver/internal/pkg/domain/model"
)

func knownAccountToDTO(tx []model.KnownAccount) []dto.KnownAccountDTO {
	knownAccountDTOList := make([]dto.KnownAccountDTO, len(tx))
	for i, element := range tx {
		accountType, err := accountTypeEnumToString(element.AccountType)
		if err != nil {
			log.Printf("Failed to convert account type to string: %v", err)
		}
		knownAccountDTOList[i] = dto.KnownAccountDTO{
			Id:            element.Id,
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
			Id:               element.Id,
			AccountNumber:    element.AccountNumber,
			AccountType:      accountType,
			AvailableBalance: element.AvailableBalance.String(),
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
			Amount:          element.Amount.String(),
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
