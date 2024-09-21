import {
  Component,
  EventEmitter,
  OnChanges,
  OnDestroy,
  OnInit,
  Optional,
  Output,
  SimpleChanges,
  SkipSelf
} from "@angular/core";
import {
  TransferFromWalletAccountDetails,
  TransferToWalletAccountDetails
} from "../../../../models/transfer-wallet-account-details";
import {Account, AccountService, KnownAccount} from "../../../../services/account/account.service";
import {TransferState, TransferType} from "../../../../models/transfer-state";
import {MatButtonToggle} from "@angular/material/button-toggle";
import {MatButton} from "@angular/material/button";
import {TransferService} from "../../../../services/transfer/transfer.service";
import {TransferFromComponent} from "../sub-components/transfer-from/transfer-from.component";
import {TransferToComponent} from "../sub-components/transfer-to/transfer-to.component";
import {TransferAmountComponent} from "../sub-components/transfer-amount/transfer-amount.component";
import {TransferTypeComponent} from "../sub-components/transfer-type/transfer-type.component";
import {Subject, takeUntil} from "rxjs";
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {FormatAccountDetailsPipe} from "../../../../pipes/format-account-details.pipe";
import {amountLessThanOrEqualToBalance} from "../../../../validators/transfer/amountLessThanOrEqualToBalance";
import {CancelWarningDialogComponent} from "../sub-components/transfer-cancel-warning/cancel-warning-dialog.component";
import {MatDialog} from "@angular/material/dialog";
@Component({
  selector: 'app-input-details',
  standalone: true,
  imports: [
    TransferToComponent,
    MatButtonToggle,
    MatButton,
    TransferFromComponent,
    TransferFromComponent,
    TransferToComponent,
    TransferAmountComponent,
    TransferTypeComponent
  ],
  templateUrl: './input-details.component.html',
  styleUrl: './input-details.component.css',
  providers: [FormatAccountDetailsPipe],
})
export class InputDetailsComponent implements OnInit, OnDestroy {
  protected hasSubmitted: boolean = false;
  protected toAccountCandidates: Array<TransferToWalletAccountDetails> = [];
  protected fromAccountCandidates: Array<TransferFromWalletAccountDetails> = [];
  private ngUnsubscribe = new Subject<void>();
  protected transferForm!: FormGroup;
  protected transferTypeControl!: FormControl<TransferType | undefined>;
  protected fromAccountControl!: FormControl<TransferFromWalletAccountDetails | undefined>;
  protected toAccountControl!: FormControl<TransferToWalletAccountDetails | undefined>;
  protected amountControl!: FormControl<number | undefined>;

  constructor(
    private accountService: AccountService,
    @SkipSelf() private transferService: TransferService,
    private formatAccountDetailsPipe: FormatAccountDetailsPipe,
    private dialog: MatDialog,
  ) { }
  ngOnInit() {
    this.accountService.getKnownAccounts$()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((knownAccounts) => {
        this.setStateToAccountCandidates(knownAccounts);
      });
    this.accountService.getCurrentAccountDetails$()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((currentAccountDetails) => {
        if (currentAccountDetails !== undefined) {
          this.setStateFromAccountCandidates(currentAccountDetails);
        }
      });
    this.initializeForm();
  }

  private setStateToAccountCandidates(knownAccounts: KnownAccount[]) {
    this.toAccountCandidates = knownAccounts.map ((account) => {
      return {
        id: account.id,
        accountNumber: account.accountNumber,
        recipientName: account.accountHolder,
        accountType: this.formatAccountType(account.accountType)
      }
    });
  }

  private setStateFromAccountCandidates(accounts: Account[]) {
    this.fromAccountCandidates = accounts.map ((account) => {
      return {
        id: account.id,
        accountNumber: account.accountNumber,
        accountHolder: account.accountHolder,
        accountType: this.formatAccountType(account.accountType),
        availableBalance: account.availableBalance,
      }
    });
  }

  private formatAccountType(accountType: string): string {
    return accountType.charAt(0).toUpperCase() + accountType.substring(1).toLowerCase()
  }

  protected validateAndProceed() {
    this.hasSubmitted = true;
    if (this.transferForm.valid) {
      const amount = this.amountControl.value!;
      const fromAccount = this.fromAccountControl.value!;
      const toAccount = this.toAccountControl.value!;
      const transferType = this.transferTypeControl.value!;
      const transferState: TransferState = {
        fromAccountNumber: this.formatAccountDetailsPipe.transformAccountNumberWithType(fromAccount),
        fromAccountId: fromAccount.id,
        toAccountNumber: this.formatAccountDetailsPipe.transformAccountNumberWithType(toAccount),
        toAccountId: toAccount.id,
        amount: amount,
        transferType: transferType,
        recipientName: toAccount.recipientName,
      }
      this.transferService.setTransferData(transferState);
    }
  }

  private initializeForm(): void {
    this.transferTypeControl = new FormControl<TransferType | undefined>(
      undefined, {nonNullable: true, validators: [Validators.required]}
    );

    this.fromAccountControl = new FormControl<TransferFromWalletAccountDetails | undefined>(
      undefined, {nonNullable: true, validators: [Validators.required]}
      );

    this.toAccountControl = new FormControl<TransferToWalletAccountDetails | undefined>(
      undefined, {nonNullable: true, validators: [Validators.required]}
    );

    this.amountControl = new FormControl<number | undefined>(
      undefined, {
        nonNullable: true,
        validators: [
          ...this.amountValidatorList,
          amountLessThanOrEqualToBalance(this.fromAccountControl.value?.availableBalance || 0)
        ]
      }
    );

    this.fromAccountControl.valueChanges
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((fromAccount) => {
        this.updateBalanceForSelectedAccount(fromAccount?.availableBalance || 0, this.amountControl)
      });

    this.transferForm = new FormGroup({
      transferType: this.transferTypeControl,
      fromAccount: this.fromAccountControl,
      toAccount: this.toAccountControl,
      amount: this.amountControl,
    });
  }

  private amountValidatorList = [
    Validators.required,
    Validators.pattern(/^\d+(\.\d{1,2})?$/),
    Validators.min(0.01),
  ]

  private updateBalanceForSelectedAccount(newBalance: number, amountControl: FormControl<number | undefined>): void {
      amountControl.clearValidators();
      amountControl.setValidators([
        ...this.amountValidatorList,
        amountLessThanOrEqualToBalance(newBalance)
      ]);
      amountControl.updateValueAndValidity();
  }

  protected cancelTransfer() {
    const dialogRef = this.dialog.open(CancelWarningDialogComponent, {
      width: '500px',
    });
    dialogRef.afterClosed().subscribe(result => {
      this.transferService.cancelTransfer();
    })
  }


  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
