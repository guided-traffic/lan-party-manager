-- Add banned_users table for banning players

CREATE TABLE IF NOT EXISTS banned_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    steam_id TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    reason TEXT DEFAULT '',
    banned_by TEXT NOT NULL,
    banned_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_banned_users_steam_id ON banned_users(steam_id);
