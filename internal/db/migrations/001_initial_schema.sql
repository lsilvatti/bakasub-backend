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
    
    log_retention_days INTEGER NOT NULL DEFAULT 7,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE favorite_folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE system_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT NOT NULL,       -- INFO, WARN, ERROR
    module TEXT NOT NULL,      -- 'translate', 'video', 'api'
    message TEXT NOT NULL,     -- 'Tradução iniciada', 'Falha no FFmpeg'
    details TEXT,              -- JSON com dados extras (ex: qual arquivo, tentativa n°)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_logs_module ON system_logs(module);

-- +goose Down
DROP TABLE favorite_folders;
DROP TABLE user_config;
DROP TABLE system_logs;