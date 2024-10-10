import { Injectable } from '@angular/core';
import {TransferState} from "../../models/transfer-state";
import {BehaviorSubject, Subject, throwError} from "rxjs";
import {DtoTransactionRequestDTO, TransactionsService} from "../../backend-api";

export interface TransferData {
  toAccountNumber: string;
  toAccountId: string;
  fromAccountNumber: string;
  fromAccountId: string;
  amount: string;
  recipientName: string;
}

@Injectable()
export class TransferService {
  private transferSubject = new BehaviorSubject<TransferData | undefined>(undefined);
  private transferData$ = this.transferSubject.asObservable();
  transferStatus = new BehaviorSubject<boolean | undefined>(undefined);
  transferValidated = new Subject<void>()
  transferCancelled = new Subject<void>()
  constructor(private transactionsService: TransactionsService) { }

  setTransferData(transferState: TransferState): void {
    if (!this.isTransferStateValid(transferState)) {
      console.error('Invalid transfer data.');
      return;
    }

    const transferData: TransferData = {
      toAccountNumber: transferState.toAccountNumber!,
      toAccountId: transferState.toAccountId!,
      fromAccountNumber: transferState.fromAccountNumber!,
      fromAccountId: transferState.fromAccountId!,
      amount: transferState.amount!,
      recipientName: transferState.recipientName!,
    }
    this.transferValidated.next()
    this.transferSubject.next(transferData)
  }

  cancelTransfer(): void {
    this.clearTransferData()
    this.transferCancelled.next()
  }

  getTransferData() {
    return this.transferData$
  }

  submitTransfer() {
    const transferData = this.transferSubject.value
    if (!transferData) {
      this.clearTransferData();
      return;
    }

    const transactionRequest: DtoTransactionRequestDTO = {
      amount: transferData.amount,
      fromBankAccountId: transferData.fromAccountId,
      toBankAccountId: transferData.toAccountId,
    }

    this.transactionsService.transactionsPost(transactionRequest)
      .subscribe({
        next: () => {
          this.transferStatus.next(true)
        },
        error: (error) => {
          this.transferStatus.next(false)
        }
      })
  }

  private isTransferStateValid(transferState: TransferState): boolean {
    return !!(
      transferState.toAccountNumber &&
      transferState.fromAccountNumber &&
      transferState.amount &&
      transferState.recipientName
    );
  }

  clearTransferData() {
    this.transferSubject.next(undefined)
    this.transferStatus.next(undefined)
  }
}
