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

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE favorite_folders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE system_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    level TEXT NOT NULL,
    module TEXT NOT NULL,
    message TEXT NOT NULL,
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_logs_module ON system_logs(module);

CREATE TABLE translation_memory (
    hash TEXT PRIMARY KEY,
    source_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    target_lang TEXT NOT NULL,
    preset TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_translation_memory_hash ON translation_memory(hash);

CREATE TABLE translation_jobs (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL, -- 'pending', 'analyzing', 'processing', 'completed', 'failed'
    file_path TEXT NOT NULL,
    target_lang TEXT NOT NULL,
    model TEXT NOT NULL,
    preset TEXT NOT NULL,
    
    -- Estimativas (Pre-flight)
    estimated_lines INTEGER DEFAULT 0,
    estimated_tokens INTEGER DEFAULT 0,
    estimated_cost_usd REAL DEFAULT 0.0,
    
    -- Progresso Real
    processed_lines INTEGER DEFAULT 0,
    total_tokens_used INTEGER DEFAULT 0,
    actual_cost_usd REAL DEFAULT 0.0,
    
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE favorite_folders;
DROP TABLE user_config;
DROP TABLE system_logs;
DROP TABLE translation_memory;