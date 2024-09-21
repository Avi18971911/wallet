import { Component } from '@angular/core';
import {MatDialogActions, MatDialogContent, MatDialogRef, MatDialogTitle} from '@angular/material/dialog';
import {MatIcon} from "@angular/material/icon";
import {MatButton} from "@angular/material/button";

@Component({
  selector: 'app-cancel-warning-dialog',
  templateUrl: './cancel-warning-dialog.component.html',
  styleUrls: ['./cancel-warning-dialog.component.css'],
  standalone: true,
  imports: [
    MatIcon,
    MatDialogActions,
    MatDialogContent,
    MatDialogTitle,
    MatButton
  ],
})
export class CancelWarningDialogComponent {
  constructor(public dialogRef: MatDialogRef<CancelWarningDialogComponent>) {}

  onAcceptance() {
    this.dialogRef.close();
  }
}
