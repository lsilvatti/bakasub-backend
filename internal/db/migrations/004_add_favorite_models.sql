-- +goose Up
ALTER TABLE user_config ADD COLUMN favorite_models JSONB DEFAULT '[]'::jsonb;

-- +goose Down
ALTER TABLE user_config DROP COLUMN favorite_models;