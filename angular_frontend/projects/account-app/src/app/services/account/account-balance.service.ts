import { Injectable } from '@angular/core';
import {AccountsService, DtoAccountBalanceHistoryRequestDTO} from "../../backend-api";
import {BehaviorSubject, Observable} from "rxjs";

export interface accountBalanceHistoryMonthsOutput {
    accountId: string;
    accountBalanceHistoryMonths: accountBalanceHistoryMonth[];
}

interface accountBalanceHistoryMonth {
  month: number;
  year: number;
  monthEndAvailableBalance: string;
  monthEndPendingBalance: string;
}

export interface accountBalanceHistoryMonthsInput {
    accountId: string;
    fromTime: string;
    toTime: string;
}

@Injectable({
  providedIn: 'root'
})
export class AccountBalanceService {
  private accountBalanceHistorySubject: BehaviorSubject<accountBalanceHistoryMonthsOutput | undefined> =
      new BehaviorSubject<accountBalanceHistoryMonthsOutput | undefined>(undefined);
  public accountBalanceHistory$: Observable<accountBalanceHistoryMonthsOutput | undefined> =
      this.accountBalanceHistorySubject

  constructor(private backendAccountService: AccountsService) { }

  calculateAccountBalanceHistory(input: accountBalanceHistoryMonthsInput): void {
      this.backendAccountService.accountsHistoryGet(
          {
              bankAccountId: input.accountId,
              fromTime: input.fromTime,
              toTime: input.toTime
          }
      ).subscribe((data) => {
          const output: accountBalanceHistoryMonthsOutput = {
            accountId: input.accountId,
            accountBalanceHistoryMonths: data.months.map((month) => {
                return {
                  month: month.month,
                  year: month.year,
                  monthEndAvailableBalance: month.availableBalance,
                  monthEndPendingBalance: month.pendingBalance
                }
            })
          }
          this.accountBalanceHistorySubject.next(output);
      });
  }

  getAccountBalanceHistory(): Observable<accountBalanceHistoryMonthsOutput | undefined> {
        return this.accountBalanceHistory$
  }
}
