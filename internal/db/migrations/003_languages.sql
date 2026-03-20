-- +goose Up

CREATE TABLE languages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO languages (code, name) VALUES
('en', 'English'),
('es', 'Spanish'),
('fr', 'French'),
('de', 'German'),
('zh', 'Chinese'),
('ja', 'Japanese'),
('ko', 'Korean'),
('ru', 'Russian'),
('pt', 'Portuguese'),
('br', 'Brazilian Portuguese');

-- +goose Down
DROP TABLE languages;