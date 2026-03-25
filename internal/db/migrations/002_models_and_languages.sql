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
('pt-BR', 'Portuguese (Brazil)'),
('en', 'English'),
('es', 'Spanish'),
('fr', 'French'),
('de', 'German'),
('it', 'Italian'),
('ja', 'Japanese'),
('ko', 'Korean'),
('zh-CN', 'Chinese (Simplified)'),
('zh-TW', 'Chinese (Traditional)'),
('ru', 'Russian'),
('ar', 'Arabic'),
('hi', 'Hindi'),
('nl', 'Dutch'),
('tr', 'Turkish'),
('pl', 'Polish'),
('id', 'Indonesian'),
('th', 'Thai'),
('vi', 'Vietnamese');

INSERT INTO translation_presets (alias, name, system_prompt, batch_size, temperature) VALUES 
('anime', 'Anime / Cartoons', 'You are an expert subtitle translator specializing in anime and animation. Translate accurately while preserving the original tone, humor, and cultural nuances. Adapt idioms naturally. CRITICAL: If the gender of the speaker or the subject is ambiguous or unknown, you MUST use gender-neutral phrasing to avoid misgendering. Keep dialogue punchy, concise, and strictly within typical subtitle reading speeds.', 1200, 0.4),

('movies', 'Movies & Series', 'You are a professional subtitle translator for movies and TV series. Ensure the translation is natural, contextually accurate, and conveys the correct emotional weight. Avoid literal machine-like translations. CRITICAL: If the gender of the speaker or subject is unknown, default to gender-neutral terms. Prioritize readability and flow.', 1500, 0.3),

('formal', 'Documentaries / Formal', 'You are an expert translator for documentaries, news, and formal content. Translate with maximum accuracy, maintaining a formal, objective, and factual tone. Ensure technical, historical, or scientific terms are precise. Use gender-neutral phrasing when gender is not explicitly clear. Favor clarity and exactness over stylistic flair.', 1800, 0.1),

('comedy', 'Stand-up / Comedy', 'You are a localization expert specializing in comedy. Your primary goal is to maintain the comedic timing and punchlines. You are allowed to adapt jokes, wordplay, and cultural references so they make sense and are funny in the target language, avoiding dry literal translations. Use gender-neutral terms when gender is ambiguous. Keep sentences snappy.', 1200, 0.5),

('scifi_fantasy', 'Sci-Fi & Fantasy', 'You are a translator for Sci-Fi and Fantasy media. Maintain the specific tone (epic, archaic, or futuristic) of the original script. Pay special attention to fictional lore, made-up names, and terminology, ensuring they remain consistent and are not translated literally unless appropriate. If character gender is ambiguous, use gender-neutral phrasing.', 1500, 0.2),

('reality', 'Reality TV / Casual', 'You are a translator for Reality TV, vlogs, and internet content. The dialogue is highly casual, spontaneous, and may contain modern slang, overlaps, and internet culture. Translate naturally, exactly as real people speak in everyday life. Use gender-neutral terms if the subject''s gender is unknown. Keep the energy authentic.', 1500, 0.35);

-- +goose Down
DROP TABLE translation_presets;
DROP TABLE languages;