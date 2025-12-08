import { Component, inject, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { WebSocketService } from '../../services/websocket.service';
import { NotificationService } from '../../services/notification.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterLink, RouterLinkActive],
  template: `
    <header class="header">
      <div class="header-content">
        <div class="header-left">
          <a routerLink="/" class="logo">
            <span class="logo-icon">🎮</span>
            <span class="logo-text">LAN Party</span>
          </a>

          @if (auth.isAuthenticated()) {
            <nav class="nav">
              <a routerLink="/rate" routerLinkActive="active" class="nav-link">
                <span class="nav-icon">⭐</span>
                Rate Player
              </a>
              <a routerLink="/timeline" routerLinkActive="active" class="nav-link">
                <span class="nav-icon">📜</span>
                Timeline
              </a>
              <a routerLink="/leaderboard" routerLinkActive="active" class="nav-link">
                <span class="nav-icon">🏆</span>
                Leaderboard
              </a>
            </nav>
          }
        </div>

        <div class="header-right">
          @if (auth.isAuthenticated()) {
            <div class="credits-badge">
              <span class="credits-icon">💎</span>
              <span class="credits-count">{{ auth.credits() }}</span>
            </div>

            <div class="user-menu" (click)="toggleMenu()">
              <img
                [src]="auth.user()?.avatar_small || auth.user()?.avatar_url || '/assets/default-avatar.png'"
                [alt]="auth.user()?.username"
                class="avatar"
              />
              <span class="username">{{ auth.user()?.username }}</span>
              <span class="dropdown-arrow">▼</span>

              @if (menuOpen) {
                <div class="dropdown-menu">
                  <a [href]="auth.user()?.profile_url" target="_blank" class="dropdown-item">
                    <span>🔗</span> Steam Profile
                  </a>
                  <button (click)="logout()" class="dropdown-item logout">
                    <span>🚪</span> Logout
                  </button>
                </div>
              }
            </div>

            <div class="ws-status" [class.connected]="ws.isConnected()">
              <span class="ws-dot"></span>
            </div>
          }
        </div>
      </div>
    </header>
  `,
  styles: [`
    @use 'variables' as *;

    .header {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      height: 64px;
      background: $bg-secondary;
      border-bottom: 1px solid $border-color;
      z-index: 100;
    }

    .header-content {
      max-width: 1400px;
      margin: 0 auto;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 24px;
    }

    .header-left {
      display: flex;
      align-items: center;
      gap: 32px;
    }

    .logo {
      display: flex;
      align-items: center;
      gap: 10px;
      text-decoration: none;
      color: $text-primary;

      .logo-icon {
        font-size: 24px;
      }

      .logo-text {
        font-size: 20px;
        font-weight: 700;
        background: $gradient-primary;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
    }

    .nav {
      display: flex;
      gap: 8px;
    }

    .nav-link {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      border-radius: $radius-md;
      color: $text-secondary;
      text-decoration: none;
      font-size: 14px;
      font-weight: 500;
      transition: all $transition-fast;

      &:hover {
        background: $bg-tertiary;
        color: $text-primary;
      }

      &.active {
        background: rgba($accent-primary, 0.15);
        color: $accent-primary;
      }

      .nav-icon {
        font-size: 16px;
      }
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .credits-badge {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 6px 14px;
      background: $bg-tertiary;
      border: 1px solid $border-color;
      border-radius: $radius-full;
      font-size: 14px;
      font-weight: 600;

      .credits-icon {
        font-size: 16px;
      }

      .credits-count {
        color: $accent-primary;
      }
    }

    .user-menu {
      position: relative;
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 6px 12px 6px 6px;
      background: $bg-tertiary;
      border: 1px solid $border-color;
      border-radius: $radius-full;
      cursor: pointer;
      transition: all $transition-fast;

      &:hover {
        border-color: $border-light;
      }

      .avatar {
        width: 32px;
        height: 32px;
        border-radius: 50%;
      }

      .username {
        font-size: 14px;
        font-weight: 500;
      }

      .dropdown-arrow {
        font-size: 10px;
        color: $text-muted;
      }
    }

    .dropdown-menu {
      position: absolute;
      top: calc(100% + 8px);
      right: 0;
      min-width: 180px;
      background: $bg-card;
      border: 1px solid $border-color;
      border-radius: $radius-md;
      box-shadow: $shadow-lg;
      overflow: hidden;
      animation: fadeIn 0.15s ease;
    }

    .dropdown-item {
      display: flex;
      align-items: center;
      gap: 10px;
      width: 100%;
      padding: 12px 16px;
      background: none;
      border: none;
      color: $text-primary;
      font-size: 14px;
      text-decoration: none;
      cursor: pointer;
      transition: background $transition-fast;

      &:hover {
        background: $bg-hover;
      }

      &.logout {
        color: $accent-error;
        border-top: 1px solid $border-color;
      }
    }

    .ws-status {
      .ws-dot {
        display: block;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: $accent-error;
        transition: background $transition-fast;
      }

      &.connected .ws-dot {
        background: $accent-success;
        box-shadow: 0 0 8px $accent-success;
      }
    }

    @keyframes fadeIn {
      from {
        opacity: 0;
        transform: translateY(-8px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
  `]
})
export class HeaderComponent implements OnInit, OnDestroy {
  auth = inject(AuthService);
  ws = inject(WebSocketService);
  private notifications = inject(NotificationService);
  private subscription?: Subscription;

  menuOpen = false;

  ngOnInit(): void {
    // Connect to WebSocket when authenticated
    if (this.auth.isAuthenticated()) {
      this.ws.connect();
    }

    // Listen for vote notifications
    this.subscription = this.ws.voteReceived$.subscribe((payload) => {
      this.notifications.voteReceived(
        payload.from_username,
        payload.achievement_name,
        payload.from_avatar,
        payload.is_positive
      );
      // Refresh user data to update credits
      this.auth.refreshUser();
    });
  }

  ngOnDestroy(): void {
    this.subscription?.unsubscribe();
  }

  toggleMenu(): void {
    this.menuOpen = !this.menuOpen;
  }

  logout(): void {
    this.ws.disconnect();
    this.auth.logout();
    this.menuOpen = false;
  }
}
