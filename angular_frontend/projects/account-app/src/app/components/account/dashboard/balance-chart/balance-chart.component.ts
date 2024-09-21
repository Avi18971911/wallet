import {Component, Inject, Input, OnChanges, PLATFORM_ID} from '@angular/core';
import {BaseChartDirective} from "ng2-charts";
import {isPlatformBrowser, NgIf} from "@angular/common";

@Component({
  selector: 'app-balance-chart',
  standalone: true,
  imports: [
    BaseChartDirective,
    NgIf
  ],
  templateUrl: './balance-chart.component.html',
  styleUrl: './balance-chart.component.css'
})
export class BalanceChartComponent implements OnChanges {
  @Input({required: true, transform: transformAvailableBalance}) availableBalance: number = 10000;
  constructor(@Inject(PLATFORM_ID) private platformId: Object) {}

  public barChartOptions = {
    scaleShowVerticalLines: false,
    responsive: true
  };

  protected isBrowser = isPlatformBrowser(this.platformId)
  public barChartLabels = ['Cash & Investments', 'Credit Card', 'Loans', 'Mortgage', 'Other'];
  public barChartLegend = false;
  public barChartData: any[] = [];

  ngOnChanges(): void {
    this.updateChartData();
  }

  private updateChartData(): void {
    this.barChartData = [
      {data: [this.availableBalance, 0.00, 0.00, 0.00, 0.00], label: 'Current'},
    ];
  }
}

function transformAvailableBalance(availableBalance: number | undefined): number {
  return availableBalance ?? 0;
}

