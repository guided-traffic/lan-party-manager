import { Component, inject, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HeaderComponent } from './components/header/header.component';
import { NotificationsComponent } from './components/notifications/notifications.component';
import { AuthService } from './services/auth.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, HeaderComponent, NotificationsComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App implements OnInit {
  private authService = inject(AuthService);

  get isAuthenticated(): boolean {
    return this.authService.isAuthenticated();
  }

  ngOnInit(): void {
    // AuthService constructor already handles token check and user loading
    // No need to call checkAuth() again here
    console.log('[App] ngOnInit - isAuthenticated:', this.isAuthenticated, '- hasToken:', !!this.authService.getToken());
  }
}
