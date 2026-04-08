-- +goose Up
ALTER TABLE user_config ADD COLUMN IF NOT EXISTS default_language TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE user_config DROP COLUMN IF EXISTS default_language;
