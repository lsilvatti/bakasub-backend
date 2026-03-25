-- +goose Up
CREATE TABLE translation_jobs (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL,
    file_path TEXT NOT NULL,
    target_lang TEXT NOT NULL,
    preset TEXT NOT NULL,
    model TEXT NOT NULL,
    total_lines INTEGER DEFAULT 0,
    processed_lines INTEGER DEFAULT 0,
    prompt_tokens INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    cost_usd REAL DEFAULT 0.0,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE translation_jobs;