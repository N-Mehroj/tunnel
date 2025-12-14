-- Migration: create_users_table (DOWN)
-- Created: 2024-12-14T00:00:00Z

DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
