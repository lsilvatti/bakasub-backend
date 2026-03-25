-- +goose Up
CREATE TABLE user_config (
    id SERIAL PRIMARY KEY,
    default_model TEXT NOT NULL,
    default_preset TEXT NOT NULL,
    remove_sdh_default BOOLEAN DEFAULT FALSE,
    video_timeout_minutes INTEGER DEFAULT 30,
    log_retention_days INTEGER DEFAULT 7,
    concurrent_translations INTEGER DEFAULT 5,
    max_retries INTEGER DEFAULT 3,
    base_retry_delay INTEGER DEFAULT 2
);

CREATE TABLE folders (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE
);

CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    level TEXT NOT NULL,
    module TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE translation_memory (
    hash TEXT PRIMARY KEY,
    source_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    target_lang TEXT NOT NULL,
    preset TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO user_config (default_model, default_preset) 
VALUES ('google/gemini-2.5-flash-lite-preview-09-2025', 'anime');

-- +goose Down
DROP TABLE translation_memory;
DROP TABLE logs;
DROP TABLE folders;
DROP TABLE user_config;