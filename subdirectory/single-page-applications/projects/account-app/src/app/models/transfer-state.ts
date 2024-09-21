export enum TransferType {
  IMMEDIATE = "IMMEDIATE",
  SCHEDULED = "SCHEDULED"
}

export interface TransferState {
  toAccountNumber: string | undefined;
  toAccountId: string | undefined;
  fromAccountNumber: string | undefined;
  fromAccountId: string | undefined;
  amount: number | undefined;
  transferType: TransferType | undefined;
  recipientName: string | undefined;
}
