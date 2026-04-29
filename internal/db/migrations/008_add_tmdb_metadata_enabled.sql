-- +goose Up
ALTER TABLE user_config ADD COLUMN tmdb_metadata_enabled BOOLEAN NOT NULL DEFAULT true;

-- +goose Down
ALTER TABLE user_config DROP COLUMN tmdb_metadata_enabled;
