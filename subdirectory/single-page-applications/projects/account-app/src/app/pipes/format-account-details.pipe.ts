import { Pipe, PipeTransform } from '@angular/core';
import {
  TransferFromWalletAccountDetails,
  TransferToWalletAccountDetails
} from "../models/transfer-wallet-account-details";

@Pipe({
  name: 'formatAccountDetails',
  standalone: true
})
export class FormatAccountDetailsPipe implements PipeTransform {

  transform(accountDetails: TransferToWalletAccountDetails | TransferFromWalletAccountDetails): string {
    return `Wallet ${accountDetails.accountType} Account ${accountDetails.accountNumber}
    ${"recipientName" in accountDetails ? accountDetails.recipientName : accountDetails.accountHolder}`;
  }

  transformAccountNumberWithType(
    accountDetails: TransferFromWalletAccountDetails | TransferToWalletAccountDetails
  ): string {
    return `Wallet ${accountDetails.accountType} Account ${accountDetails.accountNumber}`;
  }

}
