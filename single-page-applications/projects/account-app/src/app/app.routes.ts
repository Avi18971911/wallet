import { Routes } from '@angular/router';
import {LoginComponent} from "./components/login/login.component";
import {AccountComponent} from "./components/account/account.component";
import {DashboardComponent} from "./components/account/dashboard/dashboard.component";
import {TransferWalletBankComponent} from "./components/account/transfer/transfer-wallet-bank.component";
import {InputDetailsComponent} from "./components/account/transfer/input-details/input-details.component";
import {RouteNames} from "./route-names";

export const routes: Routes = [
  { path: "", redirectTo: "/login", pathMatch: "full" },
  { path: RouteNames.LOGIN, component: LoginComponent },
  {
    path: RouteNames.ACCOUNT,
    component: AccountComponent,
    children: [
      { path: RouteNames.DASHBOARD, component: DashboardComponent },
      {
        path: RouteNames.TRANSFER,
        children:[
          {
            path: RouteNames.OTHER_WALLETBANK,
            component: TransferWalletBankComponent,
            children: [
              {path: 'input-details', component: InputDetailsComponent}
            ],
          }
        ]
      },
    ],
  },
];
