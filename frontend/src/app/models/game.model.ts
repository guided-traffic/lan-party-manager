export interface Game {
  app_id: number;
  name: string;
  header_image_url: string;
  capsule_image_url: string;
  playtime_forever: number;
  categories: string[];
  owner_count: number;
  owners: string[];
  is_pinned: boolean;
  // Price information
  is_free: boolean;
  price_cents: number;
  original_cents: number;
  discount_percent: number;
  price_formatted: string;
  // Review information
  review_score: number; // Percentage of positive reviews (0-100), -1 if not enough reviews
}

export interface SyncStatus {
  needs_sync: boolean;
  is_syncing: boolean;
  phase: string;
  current_game: string;
  processed: number;
  total: number;
}

export interface GamesResponse {
  pinned_games: Game[];
  all_games: Game[];
  sync_status?: SyncStatus;
}

export interface RefreshMyGamesResponse {
  message: string;
  game_count: number;
  next_refresh_at?: string;
  warning?: string;
}

export interface RefreshMyGamesError {
  error: string;
  remaining_seconds?: number;
  cooldown_ends_at?: string;
}
