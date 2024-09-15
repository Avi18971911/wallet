import {Component, OnInit} from '@angular/core';
import {HeaderComponent} from "./header/header.component";
import {ActivatedRoute, Router, RouterOutlet} from "@angular/router";
import {RouteNames} from "../../route-names";
import {TransferService} from "../../services/transfer/transfer.service";
import {pipe, Subject, takeUntil} from "rxjs";
import {subscribe} from "node:diagnostics_channel";

@Component({
  selector: 'app-welcome',
  standalone: true,
  imports: [
    HeaderComponent,
    RouterOutlet
  ],
  templateUrl: './account.component.html',
  styleUrl: './account.component.css',
  providers: [TransferService],
})
export class AccountComponent implements OnInit {
  private ngUnsubscribe = new Subject<void>();
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private transferService: TransferService,
  ) {}
  ngOnInit() {
    this.navigateToDashboard()
  }

  private navigateToDashboard() {
    this.router.navigate([RouteNames.DASHBOARD], {relativeTo: this.route}).catch((error) => {
      console.error(error);
    });
  }
}
