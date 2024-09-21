import {Component, EventEmitter, OnDestroy, OnInit, Output} from '@angular/core';
import {TransferData, TransferService} from "../../../../services/transfer/transfer.service";
import {Subject, takeUntil} from "rxjs";
import {MatButton} from "@angular/material/button";
import {MatDialog} from "@angular/material/dialog";
import {CancelWarningDialogComponent} from "../sub-components/transfer-cancel-warning/cancel-warning-dialog.component";

@Component({
  selector: 'app-verify-details',
  standalone: true,
  imports: [
    MatButton,
  ],
  templateUrl: './verify-details.component.html',
  styleUrl: './verify-details.component.css'
})
export class VerifyDetailsComponent implements OnInit, OnDestroy {
  transferData: TransferData | undefined;
  private ngUnsubscribe = new Subject<void>();
  @Output() cancelTransaction = new EventEmitter<void>();
  constructor(
    private transferService: TransferService,
    private dialog: MatDialog,
  ) {}

  protected submitTransfer() {
    this.transferService.submitTransfer();
  }

  protected cancelTransfer() {
    const dialogRef = this.dialog.open(CancelWarningDialogComponent, {
      width: '500px',
    });
    dialogRef.afterClosed().subscribe(result => {
      this.transferService.cancelTransfer();
    })
  }

  ngOnInit() {
    this.transferService.getTransferData()
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((transferData) => {
        this.transferData = transferData;
      });
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
