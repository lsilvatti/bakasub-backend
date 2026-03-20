-- +goose Up

CREATE TABLE user_config (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    default_model TEXT NOT NULL DEFAULT 'google/gemini-2.5-flash-lite-preview-09-2025',
    default_preset TEXT NOT NULL DEFAULT 'anime',
    remove_sdh_default BOOLEAN NOT NULL DEFAULT 0,
    
    concurrent_translations INTEGER NOT NULL DEFAULT 5,
    max_retries INTEGER NOT NULL DEFAULT 3,
    base_retry_delay INTEGER NOT NULL DEFAULT 2,
    
    video_timeout_minutes INTEGER NOT NULL DEFAULT 20,
    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE favorite_folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE favorite_folders;
DROP TABLE user_config;