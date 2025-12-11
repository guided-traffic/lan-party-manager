-- Initial schema for rate-your-mate (SQLite)

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    steam_id TEXT UNIQUE NOT NULL,
    username TEXT NOT NULL,
    avatar_url TEXT,
    avatar_small TEXT,
    profile_url TEXT,
    credits INTEGER DEFAULT 0,
    last_credit_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create index for steam_id lookups
CREATE INDEX IF NOT EXISTS idx_users_steam_id ON users(steam_id);

-- Votes table
CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    to_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id TEXT NOT NULL,
    points INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CHECK (from_user_id != to_user_id)
);

-- Index for leaderboard queries
CREATE INDEX IF NOT EXISTS idx_votes_achievement ON votes(achievement_id, to_user_id);

-- Index for timeline queries
CREATE INDEX IF NOT EXISTS idx_votes_timeline ON votes(created_at DESC);

-- Chat messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    achievements TEXT DEFAULT '[]',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for chat timeline queries
CREATE INDEX IF NOT EXISTS idx_chat_messages_timeline ON chat_messages(created_at DESC);

-- Game cache table for Steam Store data
CREATE TABLE IF NOT EXISTS game_cache (
    app_id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    categories TEXT DEFAULT '[]',
    is_free INTEGER DEFAULT 0,
    price_cents INTEGER DEFAULT 0,
    original_cents INTEGER DEFAULT 0,
    discount_percent INTEGER DEFAULT 0,
    price_formatted TEXT DEFAULT '',
    fetch_failed INTEGER DEFAULT 0,
    review_score INTEGER DEFAULT -1,
    fetched_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for stale game lookups
CREATE INDEX IF NOT EXISTS idx_game_cache_fetched ON game_cache(fetched_at);
