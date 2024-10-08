package handlers

import (
	"webserver/internal/app/server/dto"
	"webserver/internal/pkg/domain/model"
)

func knownAccountToDTO(tx []model.KnownBankAccount) []dto.KnownBankAccountDTO {
	knownAccountDTOList := make([]dto.KnownBankAccountDTO, len(tx))
	for i, element := range tx {
		accountType := string(element.AccountType)
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
		accountDTOList[i] = dto.BankAccountDTO{
			Id:               element.Id,
			AccountNumber:    element.AccountNumber,
			AccountType:      element.AccountType,
			PendingBalance:   element.PendingBalance.String(),
			AvailableBalance: element.AvailableBalance.String(),
		}
	}
	return accountDTOList
}

func accountDetailsToDTO(tx *model.AccountDetailsOutput) dto.AccountDetailsResponseDTO {
	return dto.AccountDetailsResponseDTO{
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

func accountTransactionToDTO(tx []model.BankAccountTransactionOutput) []dto.AccountTransactionResponseDTO {
	accountTransactionDTOList := make([]dto.AccountTransactionResponseDTO, len(tx))
	for i, element := range tx {
		accountTransactionDTOList[i] = dto.AccountTransactionResponseDTO{
			Id:                 element.Id,
			BankAccountId:      element.BankAccountId,
			OtherBankAccountId: element.OtherBankAccountId,
			TransactionNature:  element.TransactionNature,
			TransactionType:    element.TransactionType,
			ExpirationDate:     element.ExpirationDate,
			Status:             element.Status,
			Amount:             element.Amount.String(),
			CreatedAt:          element.CreatedAt,
		}
	}
	return accountTransactionDTOList
}

func accountTransactionRequestToInput(tx *dto.AccountTransactionRequestDTO) model.TransactionsForBankAccountInput {
	return model.TransactionsForBankAccountInput{
		BankAccountId: tx.BankAccountId,
		FromTime:      tx.FromTime,
		ToTime:        tx.ToTime,
	}
}

func accountHistoryRequestToInput(tx *dto.AccountBalanceHistoryRequestDTO) model.AccountHistoryInMonthsInput {
	return model.AccountHistoryInMonthsInput{
		BankAccountId: tx.BankAccountId,
		FromTime:      tx.FromTime,
		ToTime:        tx.ToTime,
	}
}

func accountHistoryToDTO(tx *model.AccountBalanceMonthsOutput) dto.AccountBalanceHistoryResponseDTO {
	months := make([]dto.AccountBalanceMonthDTO, len(tx.Months))
	for i, element := range tx.Months {
		months[i] = dto.AccountBalanceMonthDTO{
			Month:            int(element.Month),
			Year:             element.Year,
			AvailableBalance: element.AvailableBalance.String(),
			PendingBalance:   element.PendingBalance.String(),
		}
	}
	return dto.AccountBalanceHistoryResponseDTO{
		BankAccountId: tx.BankAccountId,
		Months:        months,
	}
}
