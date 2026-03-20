-- +goose Up
CREATE TABLE translation_memory (
    hash TEXT PRIMARY KEY,
    source_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    target_lang TEXT NOT NULL,
    preset TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index para deixar a busca ultrarrápida mesmo com milhões de falas salvas
CREATE INDEX idx_translation_memory_hash ON translation_memory(hash);

-- +goose Down
DROP TABLE translation_memory;