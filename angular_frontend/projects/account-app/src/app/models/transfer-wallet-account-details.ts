export interface TransferToWalletAccountDetails {
  id: string;
  accountNumber: string;
  recipientName: string;
  accountType: string;
}

export interface TransferFromWalletAccountDetails {
  id: string;
  accountNumber: string;
  accountHolder: string;
  accountType: string;
  availableBalance: string;
}
