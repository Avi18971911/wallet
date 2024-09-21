import { Injectable } from '@angular/core';
import {AccountsService, DtoAccountDetailsDTO, DtoAccountDTO, DtoKnownAccountDTO, DtoPersonDTO} from "../../backend-api";
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
  availableBalance: number;
}

export interface AccountDetails {
  accountHolderFirstName: string;
  accountHolderLastName: string;
  knownAccounts: KnownAccount[];
  accounts: Account[];
}

@Injectable({
  providedIn: 'root'
})
export class AccountService {
  private userDataSubject: BehaviorSubject<AccountDetails | undefined> =
    new BehaviorSubject<AccountDetails | undefined>(undefined);
  public userData$: Observable<AccountDetails | undefined> = this.userDataSubject

  constructor(private backendAccountService: AccountsService) { }

  setUserData(data: DtoAccountDetailsDTO): void {
    const accountDetails = this.fromDtoAccountDetailsDTO(data)
    this.userDataSubject.next(accountDetails)
  }

  refreshUserData(): void {
    if (!this.userDataSubject.value) {
      return;
    }
    this.backendAccountService.accountsAccountIdGet(this.userDataSubject.value?.accounts[0].id)
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

  getFirstName$(): Observable<string | undefined> {
    return this.userData$.pipe(
      map((userData) => userData?.accountHolderFirstName)
    )
  }

  getCurrentBalance$(): Observable<number | undefined> {
    return this.userData$.pipe(
      map((userData) =>
        userData?.accounts.reduce(
          (acc, account) => acc + account.availableBalance, 0
        )
      )
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

  private fromDtoAccountDetailsDTO(dtoAccountDetailsDTO: DtoAccountDetailsDTO): AccountDetails {
    return {
      accountHolderFirstName: dtoAccountDetailsDTO.person.firstName,
      accountHolderLastName: dtoAccountDetailsDTO.person.lastName,
      knownAccounts: dtoAccountDetailsDTO.knownAccounts.map(this.fromDtoKnownAccountDTO),
      accounts: dtoAccountDetailsDTO.accounts.map(
        (dtoAccountDTO) => this.fromDtoAccountDTO(dtoAccountDTO, dtoAccountDetailsDTO.person)
      )
    }
  }

  private fromDtoKnownAccountDTO(dtoKnownAccountDTO: DtoKnownAccountDTO): KnownAccount {
    return {
      id: dtoKnownAccountDTO.id,
      accountNumber: dtoKnownAccountDTO.accountNumber,
      accountHolder: dtoKnownAccountDTO.accountHolder,
      accountType: dtoKnownAccountDTO.accountType,
    }
  }

  private fromDtoAccountDTO(dtoAccountDTO: DtoAccountDTO, dtoPersonDTO: DtoPersonDTO): Account {
    return {
      id: dtoAccountDTO.id,
      accountNumber: dtoAccountDTO.accountNumber,
      accountHolder: dtoPersonDTO.firstName + ' ' + dtoPersonDTO.lastName,
      accountType: dtoAccountDTO.accountType,
      availableBalance: dtoAccountDTO.availableBalance,
    }
  }
}
