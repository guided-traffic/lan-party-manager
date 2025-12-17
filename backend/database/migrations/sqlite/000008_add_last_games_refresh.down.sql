-- Remove last_games_refresh_at column from users table (SQLite)
-- Note: SQLite doesn't support DROP COLUMN in older versions, this is a workaround

-- Create a new table without the column
CREATE TABLE users_backup (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    steam_id TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    avatar_url TEXT,
    avatar_small TEXT,
    profile_url TEXT,
    credits INTEGER DEFAULT 0,
    last_credit_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_banned INTEGER DEFAULT 0,
    ban_reason TEXT DEFAULT ''
);

-- Copy data
INSERT INTO users_backup SELECT id, steam_id, username, avatar_url, avatar_small, profile_url, credits, last_credit_at, created_at, updated_at, is_banned, ban_reason FROM users;

-- Drop old table
DROP TABLE users;

-- Rename backup to original
ALTER TABLE users_backup RENAME TO users;

-- Recreate index
CREATE INDEX IF NOT EXISTS idx_users_steam_id ON users(steam_id);
