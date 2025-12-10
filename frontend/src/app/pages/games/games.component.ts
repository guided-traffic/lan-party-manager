import { Component, OnInit, signal, inject, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { GameService } from '../../services/game.service';
import { UserService } from '../../services/user.service';
import { AuthService } from '../../services/auth.service';
import { Game } from '../../models/game.model';
import { User } from '../../models/user.model';

@Component({
  selector: 'app-games',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="games-page">
      <div class="page-header">
        <h1>
          <span class="page-icon">🎮</span>
          Multiplayer Games
        </h1>
        <p class="subtitle">Spiele die von LAN-Party Teilnehmern besessen werden</p>
        <div class="header-buttons">
          <button class="refresh-btn" (click)="refreshGames()" [disabled]="loading()">
            <span class="refresh-icon" [class.spinning]="loading()">🔄</span>
            Aktualisieren
          </button>
          @if (isAdmin()) {
            <button class="refresh-btn admin-btn" (click)="invalidateCache()" [disabled]="loading() || invalidating()">
              <span class="refresh-icon" [class.spinning]="invalidating()">☁️</span>
              Update von Steam
            </button>
          }
        </div>
      </div>

      @if (loading()) {
        <div class="loading">
          <div class="spinner"></div>
          <p>Lade Spiele von allen Teilnehmern...</p>
        </div>
      } @else if (error()) {
        <div class="error">
          <span class="error-icon">❌</span>
          <p>{{ error() }}</p>
          <button (click)="loadGames()">Erneut versuchen</button>
        </div>
      } @else {
        <!-- Pinned Games Section -->
        @if (pinnedGames().length > 0) {
          <section class="games-section pinned-section">
            <h2>
              <span class="section-icon">📌</span>
              Pinned Games
              <span class="count">{{ pinnedGames().length }}</span>
            </h2>
            <div class="games-grid">
              @for (game of pinnedGames(); track game.app_id) {
                <div class="game-card pinned" (click)="openSteamStore(game.app_id)">
                  <div class="game-image">
                    <img [src]="game.header_image_url" [alt]="game.name" loading="lazy" />
                  </div>
                  <div class="game-info">
                    <div class="game-title-row">
                      <h3>{{ game.name }}</h3>
                      @if (game.price_formatted) {
                        <div class="price-badge" [class.free]="game.is_free" [class.discount]="game.discount_percent > 0">
                          @if (game.discount_percent > 0) {
                            <span class="discount-tag">-{{ game.discount_percent }}%</span>
                          }
                          <span class="price">{{ game.price_formatted }}</span>
                        </div>
                      }
                    </div>
                    <div class="game-meta">
                      @if (game.owner_count > 0) {
                        <div class="owners" [title]="getOwnerNames(game.owners)">
                          <span class="owner-icon">👥</span>
                          <span>{{ game.owner_count }} {{ game.owner_count === 1 ? 'Besitzer' : 'Besitzer' }}</span>
                        </div>
                      } @else {
                        <div class="owners no-owners">
                          <span class="owner-icon">👤</span>
                          <span>Kein Teilnehmer besitzt dieses Spiel</span>
                        </div>
                      }
                      <div class="categories">
                        @for (cat of getMultiplayerCategories(game.categories); track cat) {
                          <span class="category-tag">{{ cat }}</span>
                        }
                      </div>
                    </div>
                  </div>
                </div>
              }
            </div>
          </section>
        }

        <!-- All Games Section -->
        <section class="games-section">
          <h2>
            <span class="section-icon">🎲</span>
            Alle Multiplayer Spiele
            <span class="count">{{ allGames().length }}</span>
          </h2>

          @if (allGames().length === 0 && pinnedGames().length === 0) {
            <div class="empty-state">
              <span class="empty-icon">🎮</span>
              <p>Noch keine Multiplayer-Spiele gefunden.</p>
              <p class="hint">Spiele werden geladen sobald Spieler sich anmelden.</p>
            </div>
          } @else {
            <div class="games-grid">
              @for (game of allGames(); track game.app_id) {
                <div class="game-card" (click)="openSteamStore(game.app_id)">
                  <div class="game-image">
                    <img [src]="game.header_image_url" [alt]="game.name" loading="lazy" />
                    <div class="owner-badge" [class.highlight]="game.owner_count >= 3">
                      {{ game.owner_count }}x
                    </div>
                  </div>
                  <div class="game-info">
                    <div class="game-title-row">
                      <h3>{{ game.name }}</h3>
                      @if (game.price_formatted) {
                        <div class="price-badge" [class.free]="game.is_free" [class.discount]="game.discount_percent > 0">
                          @if (game.discount_percent > 0) {
                            <span class="discount-tag">-{{ game.discount_percent }}%</span>
                          }
                          <span class="price">{{ game.price_formatted }}</span>
                        </div>
                      }
                    </div>
                    <div class="game-meta">
                      <div class="owners" [title]="getOwnerNames(game.owners)">
                        <span class="owner-icon">👥</span>
                        <span>{{ game.owner_count }} {{ game.owner_count === 1 ? 'Besitzer' : 'Besitzer' }}</span>
                      </div>
                      <div class="categories">
                        @for (cat of getMultiplayerCategories(game.categories); track cat) {
                          <span class="category-tag">{{ cat }}</span>
                        }
                      </div>
                    </div>
                  </div>
                </div>
              }
            </div>
          }
        </section>
      }
    </div>
  `,
  styles: [`
    @use 'variables' as *;

    .games-page {
      max-width: 1200px;
      margin: 0 auto;
      padding: 24px;
    }

    .page-header {
      text-align: center;
      margin-bottom: 32px;
      position: relative;

      h1 {
        font-size: 2rem;
        margin-bottom: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
      }

      .page-icon {
        font-size: 2.5rem;
      }

      .subtitle {
        color: $text-secondary;
        font-size: 1rem;
        margin-bottom: 16px;
      }

      .header-buttons {
        display: flex;
        gap: 12px;
        justify-content: center;
        flex-wrap: wrap;
      }

      .refresh-btn {
        background: $bg-tertiary;
        border: 1px solid $border-color;
        color: $text-primary;
        padding: 8px 16px;
        border-radius: 8px;
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        gap: 8px;
        transition: all 0.2s;

        &:hover:not(:disabled) {
          background: $bg-secondary;
          border-color: $accent-primary;
        }

        &:disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }

        &.admin-btn {
          border-color: $accent-warning;

          &:hover:not(:disabled) {
            background: rgba($accent-warning, 0.1);
            border-color: $accent-warning;
          }
        }

        .refresh-icon {
          display: inline-block;
          transition: transform 0.3s;

          &.spinning {
            animation: spin 1s linear infinite;
          }
        }
      }
    }

    @keyframes spin {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }

    .loading {
      text-align: center;
      padding: 60px 20px;

      .spinner {
        width: 50px;
        height: 50px;
        border: 3px solid $border-color;
        border-top-color: $accent-primary;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin: 0 auto 16px;
      }

      p {
        color: $text-secondary;
      }
    }

    .error {
      text-align: center;
      padding: 60px 20px;
      background: rgba($accent-error, 0.1);
      border-radius: 12px;
      border: 1px solid rgba($accent-error, 0.3);

      .error-icon {
        font-size: 3rem;
        display: block;
        margin-bottom: 16px;
      }

      p {
        color: $text-primary;
        margin-bottom: 16px;
      }

      button {
        background: $accent-primary;
        color: white;
        border: none;
        padding: 10px 20px;
        border-radius: 8px;
        cursor: pointer;

        &:hover {
          opacity: 0.9;
        }
      }
    }

    .games-section {
      margin-bottom: 40px;

      h2 {
        display: flex;
        align-items: center;
        gap: 10px;
        font-size: 1.5rem;
        margin-bottom: 20px;
        padding-bottom: 12px;
        border-bottom: 2px solid $border-color;

        .section-icon {
          font-size: 1.5rem;
        }

        .count {
          background: $bg-tertiary;
          color: $text-secondary;
          font-size: 0.875rem;
          padding: 4px 10px;
          border-radius: 12px;
          font-weight: normal;
        }
      }

      &.pinned-section h2 {
        border-bottom-color: $accent-primary;
      }
    }

    .games-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 20px;
    }

    .game-card {
      background: $bg-secondary;
      border-radius: 12px;
      overflow: hidden;
      border: 1px solid $border-color;
      transition: all 0.2s;
      cursor: pointer;

      &:hover {
        transform: translateY(-4px);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
        border-color: $accent-primary;
      }

      &.pinned {
        border-color: $accent-primary;
        box-shadow: 0 0 20px rgba($accent-primary, 0.2);
      }

      .game-image {
        position: relative;
        aspect-ratio: 460 / 215;
        overflow: hidden;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.3s;
        }

        &:hover img {
          transform: scale(1.05);
        }

        .pinned-badge {
          position: absolute;
          top: 8px;
          left: 8px;
          background: rgba($accent-primary, 0.9);
          padding: 4px 8px;
          border-radius: 6px;
          font-size: 1rem;
        }

        .owner-badge {
          position: absolute;
          top: 8px;
          right: 8px;
          background: rgba(0, 0, 0, 0.8);
          color: $text-primary;
          padding: 4px 10px;
          border-radius: 6px;
          font-size: 0.875rem;
          font-weight: 600;

          &.highlight {
            background: $accent-success;
            color: white;
          }
        }

      }

      .game-info {
        padding: 16px;

        .game-title-row {
          display: flex;
          align-items: flex-start;
          justify-content: space-between;
          gap: 12px;
          margin-bottom: 10px;

          h3 {
            font-size: 1rem;
            line-height: 1.3;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
            overflow: hidden;
            margin: 0;
            flex: 1;
          }

          .price-badge {
            background: $bg-tertiary;
            color: $text-primary;
            padding: 4px 10px;
            border-radius: 6px;
            font-size: 0.875rem;
            font-weight: 600;
            display: flex;
            align-items: center;
            gap: 6px;
            flex-shrink: 0;
            white-space: nowrap;

            &.free {
              background: rgba($accent-success, 0.9);
              color: white;
            }

            &.discount {
              .discount-tag {
                background: $accent-success;
                color: white;
                padding: 2px 6px;
                border-radius: 4px;
                font-size: 0.75rem;
                font-weight: 700;
              }
            }
          }
        }

        .game-meta {
          display: flex;
          flex-direction: column;
          gap: 8px;
        }

        .owners {
          display: flex;
          align-items: center;
          gap: 6px;
          color: $text-secondary;
          font-size: 0.875rem;

          .owner-icon {
            font-size: 1rem;
          }

          &.no-owners {
            color: $accent-error;
            opacity: 0.8;
          }
        }

        .categories {
          display: flex;
          flex-wrap: wrap;
          gap: 6px;
        }

        .category-tag {
          background: $bg-tertiary;
          color: $text-secondary;
          padding: 2px 8px;
          border-radius: 4px;
          font-size: 0.75rem;
        }
      }
    }

    .empty-state {
      text-align: center;
      padding: 60px 20px;
      background: $bg-secondary;
      border-radius: 12px;
      border: 1px dashed $border-color;

      .empty-icon {
        font-size: 4rem;
        display: block;
        margin-bottom: 16px;
        opacity: 0.5;
      }

      p {
        color: $text-secondary;
        margin-bottom: 8px;
      }

      .hint {
        font-size: 0.875rem;
        opacity: 0.7;
      }
    }

    // Responsive
    @media (max-width: 768px) {
      .games-page {
        padding: 16px;
      }

      .page-header h1 {
        font-size: 1.5rem;
      }

      .games-grid {
        grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
        gap: 16px;
      }
    }
  `]
})
export class GamesComponent implements OnInit {
  private gameService = inject(GameService);
  private userService = inject(UserService);
  private authService = inject(AuthService);

  loading = signal(true);
  invalidating = signal(false);
  error = signal<string | null>(null);
  pinnedGames = signal<Game[]>([]);
  allGames = signal<Game[]>([]);
  users = signal<User[]>([]);

  isAdmin = computed(() => this.authService.user()?.is_admin ?? false);

  // Map of steamId -> username for displaying owner names
  private userMap = new Map<string, string>();

  ngOnInit() {
    this.loadUsers();
    this.loadGames();
  }

  loadUsers() {
    this.userService.getAll().subscribe({
      next: (users) => {
        this.users.set(users);
        users.forEach(u => this.userMap.set(u.steam_id, u.username));
      },
      error: (err) => console.error('Failed to load users', err)
    });
  }

  invalidateCache() {
    this.invalidating.set(true);
    this.error.set(null);

    this.gameService.invalidateCache().subscribe({
      next: () => {
        this.invalidating.set(false);
        // Now refresh to get fresh data from Steam
        this.refreshGames();
      },
      error: (err) => {
        console.error('Failed to invalidate cache', err);
        this.error.set('Fehler beim Invalidieren des Caches.');
        this.invalidating.set(false);
      }
    });
  }

  loadGames() {
    this.loading.set(true);
    this.error.set(null);

    this.gameService.getMultiplayerGames().subscribe({
      next: (response) => {
        this.pinnedGames.set(response.pinned_games || []);
        this.allGames.set(response.all_games || []);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to load games', err);
        this.error.set('Fehler beim Laden der Spiele. Bitte versuche es erneut.');
        this.loading.set(false);
      }
    });
  }

  refreshGames() {
    this.loading.set(true);
    this.error.set(null);

    this.gameService.refreshGames().subscribe({
      next: (response) => {
        this.pinnedGames.set(response.pinned_games || []);
        this.allGames.set(response.all_games || []);
        this.loading.set(false);
      },
      error: (err) => {
        console.error('Failed to refresh games', err);
        this.error.set('Fehler beim Aktualisieren der Spiele.');
        this.loading.set(false);
      }
    });
  }

  getOwnerNames(owners: string[]): string {
    if (!owners || owners.length === 0) return 'Keine Besitzer';
    return owners
      .map(steamId => this.userMap.get(steamId) || steamId)
      .join(', ');
  }

  getMultiplayerCategories(categories: string[]): string[] {
    const mpCategories = ['Multi-player', 'Co-op', 'Online Co-op', 'LAN Co-op', 'LAN PvP', 'Online PvP', 'PvP'];
    return (categories || []).filter(cat => mpCategories.includes(cat));
  }

  openSteamStore(appId: number) {
    window.open(`https://store.steampowered.com/app/${appId}`, '_blank');
  }
}
