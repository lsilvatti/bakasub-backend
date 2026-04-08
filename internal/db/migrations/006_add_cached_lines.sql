-- +goose Up
ALTER TABLE translation_jobs ADD COLUMN cached_lines INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE translation_jobs DROP COLUMN cached_lines;
