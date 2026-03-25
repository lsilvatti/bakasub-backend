-- +goose Up
CREATE TABLE languages (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE translation_presets (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    system_prompt TEXT NOT NULL,
    batch_size INTEGER NOT NULL DEFAULT 1500,
    temperature REAL NOT NULL DEFAULT 0.3
);

INSERT INTO languages (code, name) VALUES 
('pt', 'Portuguese (Brazil)'),
('en', 'English'),
('es', 'Spanish');

INSERT INTO translation_presets (alias, name, system_prompt, batch_size, temperature) VALUES 
('anime', 'Anime / Cartoons', 'You are an expert translator specializing in anime and cartoons. Translate the subtitles accurately, keeping the original tone, humor, and cultural references where possible. Keep the dialogue natural and punchy.', 1500, 0.3),
('movies', 'Movies & Series', 'You are a professional subtitle translator for movies and TV series. Ensure the translation is natural, contextually accurate, and fits typical subtitle reading speeds. Avoid overly literal translations if they sound unnatural.', 1500, 0.2),
('formal', 'Formal / Documentaries', 'You are a professional translator for documentaries and formal content. Translate with high accuracy, maintaining a formal and objective tone. Ensure technical terms are translated correctly.', 1500, 0.1);

-- +goose Down
DROP TABLE translation_presets;
DROP TABLE languages;