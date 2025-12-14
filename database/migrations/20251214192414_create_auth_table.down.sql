-- Migration: create_auth_table (DOWN)
-- Created: 2025-12-14T19:24:14+05:00

BEGIN;

DROP INDEX IF EXISTS idx_auth_tokens_user_id;
DROP INDEX IF EXISTS idx_auth_tokens_token;
DROP TABLE IF EXISTS auth_tokens;

COMMIT;
