-- Remove last_games_refresh_at column from users table (MySQL)

ALTER TABLE users DROP COLUMN last_games_refresh_at;
