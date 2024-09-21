import { Component } from '@angular/core';
import {RouterLink, RouterLinkActive} from "@angular/router";
import {DropdownMenuComponent} from "../dropdown-menu/dropdown-menu.component";
import {MatIcon} from "@angular/material/icon";
import {MenuItem} from "../../../models/menu-item.model";
import {RouteNames} from "../../../route-names";

@Component({
  selector: 'welcome-header',
  standalone: true,
  imports: [
    RouterLink,
    RouterLinkActive,
    DropdownMenuComponent,
    MatIcon,
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.css'
})
export class HeaderComponent {
  protected accountMenuTitle: string = 'My Account';
  protected transferMenuTitle: string = 'Transfer';
  protected accountMenuItems: MenuItem[] = [
    {label: 'Home', link: [RouteNames.DASHBOARD]},
    {label: 'Deposits', link: ['deposits']},
    {label: 'Credit Cards', link: ['credit-cards']},
    {label: 'Debit Cards', link: ['debit-cards']},
    {label: 'Loans', link: ['loans']},
    {label: 'Mortgages', link: ['mortgages']},
    {label: 'Other', link: ['other']},
  ];
  protected transferMenuItems: MenuItem[] = [
    {label: 'To My Accounts', link: [RouteNames.TRANSFER, 'to-my-accounts']},
    {label: 'To Other WalletBank Accounts', link: [
      RouteNames.TRANSFER,
      RouteNames.OTHER_WALLETBANK,
      RouteNames.INPUT_DETAILS
    ]},
    {label: 'To Other Banks', link: [RouteNames.TRANSFER, 'to-other-banks']},
  ];
}
