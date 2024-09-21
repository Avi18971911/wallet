import {Component, Input} from '@angular/core';
import {MatProgressSpinner} from "@angular/material/progress-spinner";
import {NgIf} from "@angular/common";

@Component({
  selector: 'app-loading-spinner',
  standalone: true,
  imports: [
    MatProgressSpinner,
    NgIf
  ],
  templateUrl: './loading-spinner.component.html',
  styleUrl: './loading-spinner.component.css'
})
export class LoadingSpinnerComponent {
  @Input() message: string = 'Loading...';
  @Input() useSpinner: boolean = true;
  @Input() diameter: number = 40;
}
