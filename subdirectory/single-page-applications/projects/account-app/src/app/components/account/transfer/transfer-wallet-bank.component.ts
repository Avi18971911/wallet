import {Component, OnDestroy, OnInit} from '@angular/core';
import {ProgressBarComponent} from "./progress-bar/progress-bar.component";
import {ActivatedRoute, NavigationStart, Router, RouterOutlet} from "@angular/router";
import {AccountService} from "../../../services/account/account.service";
import {TransferService} from "../../../services/transfer/transfer.service";
import {Subject, takeUntil} from "rxjs";
import {DateFormatService} from "../../../services/util/date-format.service";
import {RouteNames} from "../../../route-names";
import {NavigationService} from "../../../services/account/navigation.service";

@Component({
  selector: 'app-transfer',
  standalone: true,
  imports: [
    ProgressBarComponent,
    RouterOutlet
  ],
  templateUrl: './transfer-wallet-bank.component.html',
  styleUrl: './transfer-wallet-bank.component.css',
  providers: [TransferService],
})
export class TransferWalletBankComponent implements OnInit, OnDestroy {
  constructor(
    private router: Router, private route: ActivatedRoute,
    private transferService: TransferService,
    private dateService: DateFormatService,
    private accountService: AccountService,
    private navigationService: NavigationService,
  ) {}
  private ngUnsubscribe = new Subject<void>();
  protected currentStep: number = 1;
  protected dateTime: string = "";

  ngOnInit() {
    this.navigateToInputDetails()
    this.setDateTime()

    this.transferService.transferValidated
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe(() => {
        this.navigateToVerifyDetails();
      });

    this.transferService.transferStatus
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((status) => {
        if (status != undefined) {
          this.accountService.refreshUserData();
          this.navigateToTransferComplete();
        }
      });

    this.transferService.transferCancelled
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe(() => {
        this.navigationService.navigateHome();
      });

    this.router.events
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((event  ) => {
        if (event instanceof NavigationStart) {
          const targetUrl = event.url;
          if (!targetUrl.includes('/to-other-walletbank')) {
            this.transferService.clearTransferData();
          }
        }
      });
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }

  private navigateToInputDetails() {
    this.currentStep = 1;
    this.router.navigate([RouteNames.INPUT_DETAILS], {relativeTo: this.route}).catch((error) => {
      console.error(error);
    });
  }

  private navigateToVerifyDetails() {
    this.currentStep = 2;
    this.router.navigate([RouteNames.VERIFY_DETAILS], {relativeTo: this.route}).catch((error) => {
      console.error(error);
    });
  }

  private navigateToTransferComplete() {
    this.currentStep = 3;
    this.router.navigate([RouteNames.TRANSFER_COMPLETE], {relativeTo: this.route,}).catch((error) => {
      console.error(error);
    });
  }

  private setDateTime() {
    this.dateTime = this.getDateTime()
  }

  private getDateTime(): string {
    const currentDateTime = this.dateService.getCurrentDate();
    return `${currentDateTime.day} ${currentDateTime.month} ${currentDateTime.year}
    ${currentDateTime.time} ${currentDateTime.location}`;
  }

}
