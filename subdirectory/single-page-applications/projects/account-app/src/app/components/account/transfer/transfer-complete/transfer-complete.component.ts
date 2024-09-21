import {Component, Input, OnDestroy, OnInit} from '@angular/core';
import {NgIf} from "@angular/common";
import {MatCard, MatCardContent, MatCardHeader, MatCardTitle} from "@angular/material/card";
import {ActivatedRoute, Router} from "@angular/router";
import {TransferService} from "../../../../services/transfer/transfer.service";
import {Subject, takeUntil} from "rxjs";

@Component({
  selector: 'app-transfer-complete',
  standalone: true,
  imports: [
    NgIf,
    MatCardHeader,
    MatCard,
    MatCardTitle,
    MatCardContent
  ],
  templateUrl: './transfer-complete.component.html',
  styleUrl: './transfer-complete.component.css'
})
export class TransferCompleteComponent implements OnInit, OnDestroy {
  protected isSuccess: boolean = false;
  private ngUnsubscribe = new Subject<void>();
  constructor(private transferService: TransferService) {}

  ngOnInit() {
    this.transferService.transferStatus
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((isSuccess) => {
        if (isSuccess != undefined) {
          this.isSuccess = isSuccess;
        }
      });
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
