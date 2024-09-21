import {Component, EventEmitter, Input, Output} from '@angular/core';
import {TransferState, TransferType} from "../../../../../models/transfer-state";
import {MatRadioButton, MatRadioGroup} from "@angular/material/radio";
import {FormControl, ReactiveFormsModule} from "@angular/forms";
import {MatError} from "@angular/material/form-field";
import {NgIf} from "@angular/common";

@Component({
  selector: 'app-transfer-type',
  standalone: true,
  imports: [
    MatRadioButton,
    MatRadioGroup,
    ReactiveFormsModule,
    MatError,
    NgIf
  ],
  templateUrl: './transfer-type.component.html',
  styleUrl: './transfer-type.component.css'
})
export class TransferTypeComponent {
  @Input() hasSubmitted: boolean = false;
  @Input() typeControl!: FormControl<TransferType | undefined>
  protected readonly TransferType = TransferType;
}
