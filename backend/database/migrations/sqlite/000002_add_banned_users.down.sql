-- Remove banned_users table

DROP INDEX IF EXISTS idx_banned_users_steam_id;
DROP TABLE IF EXISTS banned_users;
