package handlers

import (
	"errors"
	"log"
	"webserver/internal/app/server/dto"
	"webserver/internal/pkg/domain/model"
)

func knownAccountToDTO(tx []model.KnownBankAccount) []dto.KnownBankAccountDTO {
	knownAccountDTOList := make([]dto.KnownBankAccountDTO, len(tx))
	for i, element := range tx {
		accountType, err := accountTypeEnumToString(element.AccountType)
		if err != nil {
			log.Printf("Failed to convert account type to string: %v", err)
		}
		knownAccountDTOList[i] = dto.KnownBankAccountDTO{
			Id:            element.Id,
			AccountNumber: element.AccountNumber,
			AccountHolder: element.AccountHolder,
			AccountType:   accountType,
		}
	}
	return knownAccountDTOList
}

func accountsToDTO(tx []model.BankAccount) []dto.BankAccountDTO {
	accountDTOList := make([]dto.BankAccountDTO, len(tx))
	for i, element := range tx {
		accountType, err := accountTypeEnumToString(element.AccountType)
		if err != nil {
			log.Printf("Failed to convert account type to string: %v", err)
		}
		accountDTOList[i] = dto.BankAccountDTO{
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
		BankAccounts:      accountsToDTO(tx.BankAccounts),
		KnownBankAccounts: knownAccountToDTO(tx.KnownBankAccounts),
		CreatedAt:         tx.CreatedAt,
	}
}

func accountTransactionToDTO(tx []model.BankAccountTransaction) []dto.AccountTransactionDTO {
	accountTransactionDTOList := make([]dto.AccountTransactionDTO, len(tx))
	for i, element := range tx {
		accountTransactionDTOList[i] = dto.AccountTransactionDTO{
			Id:                 element.Id,
			BankAccountId:      element.BankAccountId,
			OtherBankAccountId: element.OtherBankAccountId,
			TransactionType:    element.TransactionType,
			Amount:             element.Amount.String(),
			CreatedAt:          element.CreatedAt,
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
