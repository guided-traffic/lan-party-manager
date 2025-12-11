-- Add banned_users table for banning players

CREATE TABLE IF NOT EXISTS banned_users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    steam_id VARCHAR(20) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    reason TEXT DEFAULT (''),
    banned_by VARCHAR(20) NOT NULL,
    banned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_banned_users_steam_id (steam_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
