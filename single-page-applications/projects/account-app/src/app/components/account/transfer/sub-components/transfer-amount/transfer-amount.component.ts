import {Component, EventEmitter, Input, OnChanges, Output} from '@angular/core';
import {FormControl, FormsModule, ReactiveFormsModule} from "@angular/forms";
import {MatFormField} from "@angular/material/form-field";
import {MatInput} from "@angular/material/input";
import {TransferState} from "../../../../../models/transfer-state";
import {NgIf} from "@angular/common";
import {MatError} from "@angular/material/select";

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
  @Input() amountControl!: FormControl<number | undefined>;
  @Output() transferStateChange = new EventEmitter<Partial<TransferState>>();

  ngOnChanges() {
    if (this.hasSubmitted) {
      this.amountControl.markAsTouched({onlySelf: true});
    }
  }
}
