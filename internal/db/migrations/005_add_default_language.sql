-- +goose Up
ALTER TABLE user_config ADD COLUMN default_language TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE user_config DROP COLUMN default_language;
