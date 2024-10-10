import { Injectable } from '@angular/core';
import {
  AccountsService,
  DtoAccountDetailsResponseDTO,
  DtoBankAccountDTO,
  DtoKnownBankAccountDTO,
  DtoPersonDTO
} from "../../backend-api";
import {BehaviorSubject, map, Observable} from "rxjs";

export interface KnownAccount {
  id: string;
  accountNumber: string;
  accountHolder: string;
  accountType: string;
}

export interface Account {
  id: string;
  accountNumber: string;
  accountHolder: string;
  accountType: string;
  availableBalance: string;
}

export interface AccountDetails {
  accountHolderFirstName: string;
  accountHolderLastName: string;
  knownAccounts: KnownAccount[];
  accounts: Account[];
}

export interface FirstAndLastName {
  firstName: string;
  lastName: string;
}

@Injectable({
  providedIn: 'root'
})
export class AccountService {
  private userDataSubject: BehaviorSubject<AccountDetails | undefined> =
    new BehaviorSubject<AccountDetails | undefined>(undefined);
  public userData$: Observable<AccountDetails | undefined> = this.userDataSubject

  constructor(private backendAccountService: AccountsService) { }

  setUserData(data: DtoAccountDetailsResponseDTO): void {
    const accountDetails = this.fromDtoAccountDetailsDTO(data)
    this.userDataSubject.next(accountDetails)
  }

  refreshUserData(): void {
    if (!this.userDataSubject.value) {
      return;
    }
    this.backendAccountService.accountsAccountIdGet(this.userDataSubject.value?.accounts[0].id!)
      .subscribe((data) => {
        this.setUserData(data);
      });

  }

  clearUserData(): void {
    this.userDataSubject.next(undefined)
  }

  getKnownAccounts$(): Observable<KnownAccount[]> {
    return this.userData$.pipe(
      map((userData) => userData?.knownAccounts ?? [])
    )
  }

  getFirstAndLastName$(): Observable<FirstAndLastName | undefined> {
    return this.userData$.pipe(
      map((userData) =>
        userData? { firstName: userData.accountHolderFirstName, lastName: userData.accountHolderLastName } : undefined
      )
    );
  }

  getCurrentBalance$(): Observable<string | undefined> {
    return this.userData$.pipe(
      map((userData) => {
        const accounts = userData?.accounts
        if (!accounts) {
          return undefined
        }
        return accounts.reduce(
          (acc, account) => acc + parseFloat(account.availableBalance), 0
        ).toFixed(2)
      })
    );
  }

  getCurrentAccountDetails$(): Observable<Account[] | undefined> {
    return this.userData$.pipe(
      map((userData) => {
        if (!userData) {
          return undefined
        }
        return userData.accounts
      })
    )
  }

  private fromDtoAccountDetailsDTO(dtoAccountDetailsDTO: DtoAccountDetailsResponseDTO): AccountDetails {
    return {
      accountHolderFirstName: dtoAccountDetailsDTO.person.firstName,
      accountHolderLastName: dtoAccountDetailsDTO.person.lastName,
      knownAccounts: dtoAccountDetailsDTO.knownBankAccounts.map(this.fromDtoKnownAccountDTO),
      accounts: dtoAccountDetailsDTO.bankAccounts.map(
        (dtoAccountDTO) => this.fromDtoAccountDTO(dtoAccountDTO, dtoAccountDetailsDTO.person)
      )
    }
  }

  private fromDtoKnownAccountDTO(dtoKnownAccountDTO: DtoKnownBankAccountDTO): KnownAccount {
    return {
      id: dtoKnownAccountDTO.id,
      accountNumber: dtoKnownAccountDTO.accountNumber,
      accountHolder: dtoKnownAccountDTO.accountHolder,
      accountType: dtoKnownAccountDTO.accountType,
    }
  }

  private fromDtoAccountDTO(dtoAccountDTO: DtoBankAccountDTO, dtoPersonDTO: DtoPersonDTO): Account {
    return {
      id: dtoAccountDTO.id,
      accountNumber: dtoAccountDTO.accountNumber,
      accountHolder: dtoPersonDTO.firstName + ' ' + dtoPersonDTO.lastName,
      accountType: dtoAccountDTO.accountType,
      availableBalance: parseFloat(dtoAccountDTO.availableBalance).toFixed(2),
    }
  }
}
