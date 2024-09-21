import {Component, EventEmitter, Input, OnChanges, Output, SimpleChanges} from '@angular/core';
import {
  MatError,
  MatFormField,
  MatHint,
  MatLabel,
  MatOption,
  MatSelect,
} from "@angular/material/select";
import {TransferToWalletAccountDetails} from "../../../../../models/transfer-wallet-account-details";
import {CurrencyPipe, NgForOf, NgIf} from "@angular/common";
import {MatInput} from "@angular/material/input";
import {FormControl, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {TransferState,} from "../../../../../models/transfer-state";
import {FormatAccountDetailsPipe} from "../../../../../pipes/format-account-details.pipe";

@Component({
  selector: 'app-transfer-to',
  standalone: true,
  imports: [
    MatSelect,
    MatOption,
    NgForOf,
    MatLabel,
    MatFormField,
    NgIf,
    MatInput,
    FormsModule,
    CurrencyPipe,
    MatHint,
    MatError,
    FormatAccountDetailsPipe,
    ReactiveFormsModule,
  ],
  templateUrl: './transfer-to.component.html',
  styleUrl: './transfer-to.component.css',
})
export class TransferToComponent implements OnChanges {
  @Input() toCandidateAccountDetails: Array<TransferToWalletAccountDetails> = [];
  @Input() hasSubmitted: boolean = false;
  @Input() toControl!: FormControl<TransferToWalletAccountDetails | undefined>;

  ngOnChanges(changes: SimpleChanges) {
    if (this.hasSubmitted) {
      this.toControl.markAsTouched({onlySelf: true});
    }
  }
}
