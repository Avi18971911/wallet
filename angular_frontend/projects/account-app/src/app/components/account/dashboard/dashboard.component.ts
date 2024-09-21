import {Component, OnDestroy, OnInit} from '@angular/core';
import {DtoAccountDetailsDTO} from "../../../backend-api";
import {CommonModule} from "@angular/common";
import {BalanceChartComponent} from "./balance-chart/balance-chart.component";
import {DateFormatService} from "../../../services/util/date-format.service";
import {AccountService} from "../../../services/account/account.service";
import {LoadingSpinnerComponent} from "../../loading-spinner/loading-spinner.component";
import {Subject, takeUntil} from "rxjs";

interface CurrentMonthAndYear {
  month: string;
  year: number;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, BalanceChartComponent, LoadingSpinnerComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit, OnDestroy {
  protected userFirstName: string | undefined;
  protected currentBalance: number | undefined;
  private ngUnsubscribe = new Subject<void>();
  constructor(
    private accountService: AccountService,
    private dateService: DateFormatService,
  ) { }

  ngOnInit() {
    this.accountService.getFirstName$()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((firstName) => { this.userFirstName = firstName })
    this.accountService.getCurrentBalance$()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((currentBalance) => this.currentBalance = currentBalance)
  }

  get currentMonthAndYear(): CurrentMonthAndYear {
    const currentDateTime = this.dateService.getCurrentDate();
    const currentMonth = currentDateTime.month;
    const currentYear = parseInt(currentDateTime.year);
    return {
      month: currentMonth,
      year: currentYear,
    }
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
