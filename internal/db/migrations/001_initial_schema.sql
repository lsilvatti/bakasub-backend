-- +goose Up
CREATE TABLE user_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    default_model TEXT NOT NULL,
    default_preset TEXT NOT NULL,
    remove_sdh_default BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE favorite_folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE user_config;
DROP TABLE favorite_folders;