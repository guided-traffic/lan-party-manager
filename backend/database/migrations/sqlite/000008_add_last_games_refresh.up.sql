-- Add last_games_refresh_at column to users table (SQLite)

ALTER TABLE users ADD COLUMN last_games_refresh_at DATETIME DEFAULT NULL;
