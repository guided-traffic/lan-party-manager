import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="login-page">
      <div class="login-container">
        <div class="login-header">
          <div class="logo">
            <span class="logo-icon">🎮</span>
            <h1 class="logo-text">LAN Party Manager</h1>
          </div>
          <p class="tagline">Bewerte deine Mitspieler mit Achievements!</p>
        </div>

        <div class="login-card">
          <h2>Willkommen!</h2>
          <p class="description">
            Melde dich mit deinem Steam-Account an, um andere Spieler zu bewerten
            und Achievements zu sammeln.
          </p>

          <button class="btn btn-steam" (click)="login()">
            <svg viewBox="0 0 24 24" class="steam-icon">
              <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12c0 4.84 3.44 8.87 8 9.8V15H8v-3h2V9.5C10 7.57 11.57 6 13.5 6H16v3h-2c-.55 0-1 .45-1 1v2h3v3h-3v6.95c5.05-.5 9-4.76 9-9.95 0-5.52-4.48-10-10-10z"/>
            </svg>
            Mit Steam anmelden
          </button>

          <div class="features">
            <div class="feature">
              <span class="feature-icon">⭐</span>
              <span>Bewerte andere Spieler</span>
            </div>
            <div class="feature">
              <span class="feature-icon">🏆</span>
              <span>Sammle Achievements</span>
            </div>
            <div class="feature">
              <span class="feature-icon">📊</span>
              <span>Steige im Leaderboard auf</span>
            </div>
          </div>
        </div>

        <div class="achievements-preview">
          <h3>Verfügbare Achievements</h3>
          <div class="achievement-tags">
            <span class="achievement-chip positive">🎯 Pro Player</span>
            <span class="achievement-chip positive">👑 Endboss</span>
            <span class="achievement-chip positive">🤝 Teamplayer</span>
            <span class="achievement-chip positive">⭐ MVP</span>
            <span class="achievement-chip negative">🐣 Noob</span>
            <span class="achievement-chip negative">⛺ Camper</span>
            <span class="achievement-chip negative">😤 Rage Quitter</span>
            <span class="achievement-chip negative">☠️ Toxic</span>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    @use 'variables' as *;

    .login-page {
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px;
      background: radial-gradient(ellipse at top, #1a1a2e 0%, $bg-primary 100%);
    }

    .login-container {
      max-width: 480px;
      width: 100%;
    }

    .login-header {
      text-align: center;
      margin-bottom: 32px;
    }

    .logo {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 12px;
      margin-bottom: 12px;

      .logo-icon {
        font-size: 48px;
      }

      .logo-text {
        font-size: 32px;
        font-weight: 700;
        background: $gradient-primary;
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
    }

    .tagline {
      color: $text-secondary;
      font-size: 18px;
    }

    .login-card {
      background: $bg-card;
      border: 1px solid $border-color;
      border-radius: $radius-xl;
      padding: 40px;
      text-align: center;

      h2 {
        font-size: 24px;
        margin-bottom: 12px;
      }

      .description {
        color: $text-secondary;
        margin-bottom: 32px;
        line-height: 1.6;
      }
    }

    .btn-steam {
      width: 100%;
      padding: 16px 24px;
      font-size: 18px;
      background: linear-gradient(135deg, #1b2838, #2a475e);
      border: none;
      border-radius: $radius-md;
      color: white;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 12px;
      transition: all $transition-fast;

      &:hover {
        background: linear-gradient(135deg, #2a475e, #3a5a7c);
        transform: translateY(-2px);
        box-shadow: $shadow-lg;
      }

      .steam-icon {
        width: 24px;
        height: 24px;
      }
    }

    .features {
      margin-top: 32px;
      display: flex;
      flex-direction: column;
      gap: 12px;
    }

    .feature {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 16px;
      background: $bg-tertiary;
      border-radius: $radius-md;
      font-size: 14px;

      .feature-icon {
        font-size: 20px;
      }
    }

    .achievements-preview {
      margin-top: 32px;
      text-align: center;

      h3 {
        font-size: 16px;
        color: $text-secondary;
        margin-bottom: 16px;
      }
    }

    .achievement-tags {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      justify-content: center;
    }

    .achievement-chip {
      padding: 6px 12px;
      border-radius: $radius-full;
      font-size: 13px;
      font-weight: 500;

      &.positive {
        background: rgba($accent-positive, 0.12);
        color: $accent-positive;
        border: 1px solid rgba($accent-positive, 0.3);
      }

      &.negative {
        background: rgba($accent-negative, 0.12);
        color: $accent-negative;
        border: 1px solid rgba($accent-negative, 0.3);
      }
    }
  `]
})
export class LoginComponent {
  private auth = inject(AuthService);

  login(): void {
    this.auth.login();
  }
}
