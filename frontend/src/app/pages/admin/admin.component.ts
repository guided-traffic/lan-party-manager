import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { SettingsService } from '../../services/settings.service';
import { AuthService } from '../../services/auth.service';
import { NotificationService } from '../../services/notification.service';

@Component({
  selector: 'app-admin',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="admin-page">
      <div class="admin-container">
        <div class="admin-header">
          <h1>⚙️ Admin Panel</h1>
          <p class="admin-subtitle">Einstellungen für das Credit System</p>
        </div>

        @if (loading()) {
          <div class="loading">
            <div class="spinner"></div>
            <span>Lade Einstellungen...</span>
          </div>
        } @else if (error()) {
          <div class="error-message">
            <span>❌</span>
            <p>{{ error() }}</p>
            <button (click)="loadSettings()" class="retry-btn">Erneut versuchen</button>
          </div>
        } @else {
          <div class="settings-card">
            <div class="setting-group">
              <label for="creditInterval">Credit Interval (Minuten)</label>
              <p class="setting-description">
                Wie viele Minuten zwischen dem Verdienen von Credits vergehen.
              </p>
              <div class="input-group">
                <input
                  type="number"
                  id="creditInterval"
                  [(ngModel)]="creditIntervalMinutes"
                  min="1"
                  max="60"
                  class="setting-input"
                />
                <span class="input-suffix">min</span>
              </div>
            </div>

            <div class="setting-group">
              <label for="creditMax">Maximale Credits</label>
              <p class="setting-description">
                Die maximale Anzahl an Credits, die ein Spieler ansammeln kann.
              </p>
              <div class="input-group">
                <input
                  type="number"
                  id="creditMax"
                  [(ngModel)]="creditMax"
                  min="1"
                  max="100"
                  class="setting-input"
                />
                <span class="input-suffix">Credits</span>
              </div>
            </div>

            <div class="actions">
              <button
                (click)="saveSettings()"
                [disabled]="saving() || !hasChanges()"
                class="save-btn"
              >
                @if (saving()) {
                  <span class="btn-spinner"></span>
                  Speichern...
                } @else {
                  💾 Einstellungen speichern
                }
              </button>
              <button
                (click)="resetToOriginal()"
                [disabled]="saving() || !hasChanges()"
                class="reset-btn"
              >
                ↩️ Zurücksetzen
              </button>
            </div>

            @if (hasChanges()) {
              <div class="changes-notice">
                <span>⚠️</span>
                <span>Du hast ungespeicherte Änderungen.</span>
              </div>
            }
          </div>

          <div class="info-card">
            <h3>ℹ️ Hinweis</h3>
            <p>
              Änderungen werden <strong>sofort live</strong> an alle verbundenen Spieler übertragen.
              Die Credits-Anzeige aller Spieler wird automatisch aktualisiert.
            </p>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    @use 'variables' as *;

    .admin-page {
      min-height: calc(100vh - 64px);
      padding: 32px 24px;
      background: $bg-primary;
    }

    .admin-container {
      max-width: 600px;
      margin: 0 auto;
    }

    .admin-header {
      text-align: center;
      margin-bottom: 32px;

      h1 {
        font-size: 28px;
        font-weight: 700;
        color: $text-primary;
        margin-bottom: 8px;
      }

      .admin-subtitle {
        color: $text-secondary;
        font-size: 16px;
      }
    }

    .loading {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 16px;
      padding: 48px;
      color: $text-secondary;
    }

    .spinner {
      width: 40px;
      height: 40px;
      border: 3px solid $border-color;
      border-top-color: $accent-primary;
      border-radius: 50%;
      animation: spin 1s linear infinite;
    }

    .error-message {
      text-align: center;
      padding: 48px;
      background: $bg-card;
      border-radius: $radius-lg;
      border: 1px solid $accent-error;

      span {
        font-size: 48px;
        display: block;
        margin-bottom: 16px;
      }

      p {
        color: $accent-error;
        margin-bottom: 16px;
      }

      .retry-btn {
        padding: 10px 20px;
        background: $accent-primary;
        color: white;
        border: none;
        border-radius: $radius-md;
        cursor: pointer;
        font-weight: 500;

        &:hover {
          background: $accent-secondary;
        }
      }
    }

    .settings-card {
      background: $bg-card;
      border: 1px solid $border-color;
      border-radius: $radius-lg;
      padding: 24px;
      margin-bottom: 24px;
    }

    .setting-group {
      margin-bottom: 24px;

      &:last-of-type {
        margin-bottom: 32px;
      }

      label {
        display: block;
        font-size: 16px;
        font-weight: 600;
        color: $text-primary;
        margin-bottom: 4px;
      }

      .setting-description {
        font-size: 14px;
        color: $text-muted;
        margin-bottom: 12px;
      }
    }

    .input-group {
      display: flex;
      align-items: center;
      gap: 8px;

      .setting-input {
        width: 120px;
        padding: 12px 16px;
        background: $bg-tertiary;
        border: 1px solid $border-color;
        border-radius: $radius-md;
        color: $text-primary;
        font-size: 18px;
        font-weight: 600;

        &:focus {
          outline: none;
          border-color: $accent-primary;
          box-shadow: 0 0 0 3px rgba($accent-primary, 0.2);
        }

        &::-webkit-inner-spin-button,
        &::-webkit-outer-spin-button {
          opacity: 1;
        }
      }

      .input-suffix {
        color: $text-secondary;
        font-size: 14px;
      }
    }

    .actions {
      display: flex;
      gap: 12px;
    }

    .save-btn {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      padding: 14px 24px;
      background: $gradient-primary;
      color: white;
      border: none;
      border-radius: $radius-md;
      font-size: 16px;
      font-weight: 600;
      cursor: pointer;
      transition: all $transition-fast;

      &:hover:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: $shadow-lg;
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    }

    .reset-btn {
      padding: 14px 20px;
      background: $bg-tertiary;
      color: $text-secondary;
      border: 1px solid $border-color;
      border-radius: $radius-md;
      font-size: 14px;
      font-weight: 500;
      cursor: pointer;
      transition: all $transition-fast;

      &:hover:not(:disabled) {
        background: $bg-hover;
        color: $text-primary;
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    }

    .btn-spinner {
      width: 16px;
      height: 16px;
      border: 2px solid rgba(white, 0.3);
      border-top-color: white;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
    }

    .changes-notice {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-top: 16px;
      padding: 12px;
      background: rgba($accent-warning, 0.1);
      border: 1px solid rgba($accent-warning, 0.3);
      border-radius: $radius-md;
      color: $accent-warning;
      font-size: 14px;
    }

    .info-card {
      background: rgba($accent-primary, 0.1);
      border: 1px solid rgba($accent-primary, 0.2);
      border-radius: $radius-lg;
      padding: 20px;

      h3 {
        font-size: 16px;
        font-weight: 600;
        color: $accent-primary;
        margin-bottom: 8px;
      }

      p {
        font-size: 14px;
        color: $text-secondary;
        line-height: 1.6;

        strong {
          color: $accent-primary;
        }
      }
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }
  `]
})
export class AdminComponent implements OnInit {
  private settingsService = inject(SettingsService);
  private authService = inject(AuthService);
  private notifications = inject(NotificationService);
  private router = inject(Router);

  loading = signal(true);
  saving = signal(false);
  error = signal<string | null>(null);

  // Form values
  creditIntervalMinutes = 10;
  creditMax = 10;

  // Original values for comparison
  private originalCreditIntervalMinutes = 10;
  private originalCreditMax = 10;

  ngOnInit(): void {
    // Check if user is admin
    const user = this.authService.user();
    if (!user?.is_admin) {
      this.router.navigate(['/timeline']);
      return;
    }

    this.loadSettings();
  }

  loadSettings(): void {
    this.loading.set(true);
    this.error.set(null);

    this.settingsService.getSettings().subscribe({
      next: (settings) => {
        this.creditIntervalMinutes = settings.credit_interval_minutes;
        this.creditMax = settings.credit_max;
        this.originalCreditIntervalMinutes = settings.credit_interval_minutes;
        this.originalCreditMax = settings.credit_max;
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load settings:', err);
        this.error.set('Einstellungen konnten nicht geladen werden.');
        this.loading.set(false);
      }
    });
  }

  saveSettings(): void {
    this.saving.set(true);

    this.settingsService.updateSettings({
      credit_interval_minutes: this.creditIntervalMinutes,
      credit_max: this.creditMax
    }).subscribe({
      next: (settings) => {
        this.originalCreditIntervalMinutes = settings.credit_interval_minutes;
        this.originalCreditMax = settings.credit_max;
        this.saving.set(false);
        this.notifications.success('✅ Gespeichert', 'Einstellungen wurden gespeichert und an alle Spieler übertragen');
      },
      error: (err) => {
        console.error('Failed to save settings:', err);
        this.saving.set(false);
        this.notifications.error('❌ Fehler', 'Einstellungen konnten nicht gespeichert werden');
      }
    });
  }

  resetToOriginal(): void {
    this.creditIntervalMinutes = this.originalCreditIntervalMinutes;
    this.creditMax = this.originalCreditMax;
  }

  hasChanges(): boolean {
    return this.creditIntervalMinutes !== this.originalCreditIntervalMinutes ||
           this.creditMax !== this.originalCreditMax;
  }
}
