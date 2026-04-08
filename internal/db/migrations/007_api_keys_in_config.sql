-- +goose Up
ALTER TABLE user_config
    ADD COLUMN IF NOT EXISTS openrouter_api_key TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS tmdb_access_token  TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS concurrent_translations INTEGER NOT NULL DEFAULT 5,
    ADD COLUMN IF NOT EXISTS max_retries         INTEGER NOT NULL DEFAULT 3,
    ADD COLUMN IF NOT EXISTS base_retry_delay    INTEGER NOT NULL DEFAULT 2;

-- +goose Down
ALTER TABLE user_config
    DROP COLUMN IF EXISTS openrouter_api_key,
    DROP COLUMN IF EXISTS tmdb_access_token,
    DROP COLUMN IF EXISTS concurrent_translations,
    DROP COLUMN IF EXISTS max_retries,
    DROP COLUMN IF EXISTS base_retry_delay;
