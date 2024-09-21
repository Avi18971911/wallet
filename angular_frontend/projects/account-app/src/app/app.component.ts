import {Component, OnDestroy, OnInit} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { LoginComponent } from "./components/login/login.component";
import { HttpClientModule } from "@angular/common/http";
import { Router } from "@angular/router";
import {BarController, BarElement, CategoryScale, Chart, Colors, Legend, LinearScale, Title, Tooltip} from "chart.js";
import {AuthService} from "./services/account/auth.service";
import {AccountService} from "./services/account/account.service";
import {RouteNames} from "./route-names";
import {Subject, takeUntil} from "rxjs";

Chart.register(BarController, BarElement, CategoryScale, Colors, LinearScale, Title, Tooltip, Legend);

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    LoginComponent,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit, OnDestroy {
  title = 'account-app';
  private ngUnsubscribe = new Subject<void>();
  constructor(
    private router: Router,
    private authService: AuthService,
    private accountService: AccountService,
  ) {}
  ngOnInit() {
    this.router.navigate([RouteNames.LOGIN]);
    this.authService.isAuthenticated$
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((isAuthenticated) => {
        if (isAuthenticated) {
          this.router.navigate([RouteNames.ACCOUNT]);
        }
        else {
          this.accountService.clearUserData()
          this.router.navigate([RouteNames.LOGIN]);
        }
    });
    this.authService.loginResponse$
      .pipe(takeUntil(this.ngUnsubscribe))
      .subscribe((response) => {
        if (response) {
          this.accountService.setUserData(response);
        }
    });
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
