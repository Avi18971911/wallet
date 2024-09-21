import { Injectable } from '@angular/core';
import {Subject} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class NavigationService {
  constructor() { }
  private navigateHomeSubject = new Subject<void>();
  navigateHome$ = this.navigateHomeSubject.asObservable();

  navigateHome() {
    this.navigateHomeSubject.next();
  }
}
