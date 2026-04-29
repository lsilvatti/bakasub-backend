-- +goose Up
ALTER TABLE user_config ADD COLUMN favorite_models TEXT NOT NULL DEFAULT '[]';

-- +goose Down
ALTER TABLE user_config DROP COLUMN favorite_models;