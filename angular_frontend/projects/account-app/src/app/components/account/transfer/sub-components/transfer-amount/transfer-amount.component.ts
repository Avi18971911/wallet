import {Component, Input, OnChanges} from '@angular/core';
import {FormControl, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {MatFormField} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {NgIf} from "@angular/common";
import {MatError} from "@angular/material/select";
import {exceedBalanceErrorMessage} from "../../../../../validators/transfer/transfer-validators";

@Component({
  selector: 'app-transfer-amount',
  standalone: true,
  imports: [
    FormsModule,
    MatFormField,
    MatInput,
    MatError,
    ReactiveFormsModule,
    NgIf,
  ],
  templateUrl: './transfer-amount.component.html',
  styleUrl: './transfer-amount.component.css'
})
export class TransferAmountComponent implements OnChanges {
  @Input() hasSubmitted: boolean = false;
  @Input() amountControl!: FormControl<string | undefined>;

  ngOnChanges() {
    if (this.hasSubmitted) {
      this.amountControl.markAsTouched({onlySelf: true});
    }
  }

  protected readonly exceedBalanceErrorMessage = exceedBalanceErrorMessage;
}
