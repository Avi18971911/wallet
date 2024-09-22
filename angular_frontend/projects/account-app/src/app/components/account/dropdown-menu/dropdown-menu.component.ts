import {Component, ElementRef, Input, Renderer2} from '@angular/core';
import {MatIcon, MatIconModule} from "@angular/material/icon";
import {MatMenu, MatMenuItem, MatMenuModule, MatMenuTrigger} from "@angular/material/menu";
import {MatButtonModule, MatIconButton} from "@angular/material/button";
import {NgClass, NgForOf} from "@angular/common";
import {MenuItem} from "../../../models/menu-item.model";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-dropdown-menu',
  standalone: true,
  imports: [
    MatButtonModule,
    MatMenuModule,
    MatIconModule,
    NgForOf,
    NgClass,
    RouterLink
  ],
  templateUrl: './dropdown-menu.component.html',
  styleUrl: './dropdown-menu.component.css'
})
export class DropdownMenuComponent {
  protected timedOutCloser: number | undefined;
  protected isHovered = false;
  @Input({required: true}) title: string = "";
  @Input({required: true}) menuItems: MenuItem[] = [];
  constructor(private renderer: Renderer2) {}

  onMouseEnter(trigger: MatMenuTrigger) {
    if (this.timedOutCloser) {
      window.clearTimeout(this.timedOutCloser);
    }
    this.isHovered = true;
    trigger.openMenu();
  }

  onMouseLeave(trigger: MatMenuTrigger) {
    this.timedOutCloser = window.setTimeout(() => {
      trigger.closeMenu();
      this.isHovered = false;
    }, 150);
  }

  focusMenuItem(menuItemButton: MatMenuItem): void {
    this.renderer.selectRootElement(menuItemButton).focus();
  }

  navigationCallback(): void {
    this.isHovered = false;
    this.removeFocus();
  }

  removeFocus(): void {
    const activeElement = document.activeElement as HTMLElement;
    if (activeElement) {
      activeElement.blur();
    }
  }
}
