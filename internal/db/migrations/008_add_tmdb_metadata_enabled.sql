-- +goose Up
ALTER TABLE user_config
    ADD COLUMN IF NOT EXISTS tmdb_metadata_enabled BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE user_config
    DROP COLUMN IF EXISTS tmdb_metadata_enabled;
