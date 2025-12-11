-- Initial schema for rate-your-mate (MySQL)

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    steam_id VARCHAR(20) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    avatar_small TEXT,
    profile_url TEXT,
    credits INT DEFAULT 0,
    last_credit_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_steam_id (steam_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Votes table
CREATE TABLE IF NOT EXISTS votes (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    from_user_id BIGINT UNSIGNED NOT NULL,
    to_user_id BIGINT UNSIGNED NOT NULL,
    achievement_id VARCHAR(50) NOT NULL,
    points INT DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_votes_achievement (achievement_id, to_user_id),
    INDEX idx_votes_timeline (created_at DESC),
    CONSTRAINT chk_no_self_vote CHECK (from_user_id != to_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Chat messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    message TEXT NOT NULL,
    achievements TEXT DEFAULT ('[]'),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_chat_messages_timeline (created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Game cache table for Steam Store data
CREATE TABLE IF NOT EXISTS game_cache (
    app_id BIGINT UNSIGNED PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    categories TEXT DEFAULT ('[]'),
    is_free TINYINT(1) DEFAULT 0,
    price_cents INT DEFAULT 0,
    original_cents INT DEFAULT 0,
    discount_percent INT DEFAULT 0,
    price_formatted VARCHAR(50) DEFAULT '',
    fetch_failed TINYINT(1) DEFAULT 0,
    review_score INT DEFAULT -1,
    fetched_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_game_cache_fetched (fetched_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
