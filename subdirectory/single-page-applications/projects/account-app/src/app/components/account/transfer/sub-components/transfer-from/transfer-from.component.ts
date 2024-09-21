import {Component, Input, OnChanges, SimpleChanges} from '@angular/core';
import {MatFormField} from "@angular/material/form-field";
import {MatError, MatHint, MatOption, MatSelect} from "@angular/material/select";
import {TransferFromWalletAccountDetails} from "../../../../../models/transfer-wallet-account-details";
import {FormatAccountDetailsPipe} from "../../../../../pipes/format-account-details.pipe";
import {CurrencyPipe, NgForOf, NgIf} from "@angular/common";
import {FormControl, ReactiveFormsModule} from "@angular/forms";

@Component({
  selector: 'app-transfer-from',
  standalone: true,
  imports: [
    MatFormField,
    MatSelect,
    MatHint,
    MatOption,
    FormatAccountDetailsPipe,
    NgForOf,
    ReactiveFormsModule,
    MatError,
    NgIf,
    CurrencyPipe
  ],
  templateUrl: './transfer-from.component.html',
  styleUrl: './transfer-from.component.css',
  providers: [FormatAccountDetailsPipe],
})
export class TransferFromComponent implements OnChanges {
  constructor(private formatAccountDetailsPipe: FormatAccountDetailsPipe) {}
  protected defaultFromAccount: string = "Please select an account";

  @Input() fromCandidateAccountDetails: Array<TransferFromWalletAccountDetails> = [];
  @Input() hasSubmitted: boolean = false;
  @Input() fromControl!: FormControl<TransferFromWalletAccountDetails | undefined>;

  private getDefaultFromAccount(): string {
    if (this.fromCandidateAccountDetails.length > 0) {
      return this.formatAccountDetailsPipe.transform(this.fromCandidateAccountDetails[0]);
    }
    return "Please select an account";
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['fromControl'] && this.fromControl) {
      const defaultFromAccount = this.getDefaultFromAccount();
      if (defaultFromAccount !== "Please select an account") {
        this.fromControl.setValue(this.fromCandidateAccountDetails[0]);
      }
      this.defaultFromAccount = defaultFromAccount;
    }
  }
}
