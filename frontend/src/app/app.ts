import { Component, inject, OnInit, OnDestroy, effect } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HeaderComponent } from './components/header/header.component';
import { NotificationsComponent } from './components/notifications/notifications.component';
import { AuthService } from './services/auth.service';
import { WebSocketService } from './services/websocket.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, HeaderComponent, NotificationsComponent],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App implements OnInit, OnDestroy {
  private authService = inject(AuthService);
  private wsService = inject(WebSocketService);

  get isAuthenticated(): boolean {
    return this.authService.isAuthenticated();
  }

  constructor() {
    // Connect/disconnect WebSocket based on authentication state
    effect(() => {
      const isAuth = this.authService.isAuthenticated();
      const user = this.authService.user();

      if (isAuth && user) {
        this.wsService.connect();
      } else if (!isAuth) {
        this.wsService.disconnect();
      }
    });
  }

  ngOnInit(): void {
    console.log('[App] ngOnInit - isAuthenticated:', this.isAuthenticated, '- hasToken:', !!this.authService.getToken());
  }

  ngOnDestroy(): void {
    this.wsService.disconnect();
  }
}
