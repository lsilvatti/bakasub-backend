-- +goose Up
ALTER TABLE user_config ADD COLUMN openrouter_api_key TEXT NOT NULL DEFAULT '';
ALTER TABLE user_config ADD COLUMN tmdb_access_token TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE user_config DROP COLUMN openrouter_api_key;
ALTER TABLE user_config DROP COLUMN tmdb_access_token;
