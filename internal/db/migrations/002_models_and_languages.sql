-- +goose Up

CREATE TABLE translation_presets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    system_prompt TEXT NOT NULL,
    batch_size INTEGER NOT NULL,
    temperature REAL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO translation_presets (alias, name, system_prompt, batch_size, temperature) VALUES
(
    'movie', 
    'Movies & Series', 
    "You are a professional subtitle translator for Movies and TV Series. Translate the following dialogue into the target language. Maintain natural flow, local idioms, and emotional subtext. Ensure character gender consistency and keep sentences concise for screen reading. IMPORTANT: The text contains multiple subtitles separated by '---NEXT---'. You MUST preserve this exact separator between each translated segment. OUTPUT ONLY THE TRANSLATED TEXT. Do not add introductions, commentary, or quotes.",
    2000,
    0.5
),
(
    'anime', 
    'Anime & Animation', 
    "You are an expert anime fansubber. Translate the following dialogue into the target language. Respect honorifics (e.g., -san, -senpai) where appropriate and maintain the energetic or stylized tone typical of Japanese animation. Translate catchy phrases naturally while keeping the original 'spirit'. IMPORTANT: Use '---NEXT---' to separate each translated block. OUTPUT ONLY THE TRANSLATED TEXT. Do not merge lines or add translator notes.",
    1800,
    0.6
),
(
    'documentary', 
    'Documentary & Educational', 
    "You are a professional narrator and translator for documentaries. Use a formal, clear, and informative tone. Ensure technical terms, scientific concepts, and historical facts are translated accurately. Avoid overly colloquial language. IMPORTANT: Each subtitle is separated by '---NEXT---'. Your output must maintain this separator between every translated segment. OUTPUT ONLY THE TRANSLATED TEXT. No extra conversational text.",
    2500,
    0.2
),
(
    'youtube', 
    'YouTube & Vlogs', 
    "You are a creative translator for YouTube content. The tone should be engaging, modern, and suitable for internet culture. Translate clickbait-style phrases or energetic speech with high impact. Use common internet slang and local expressions if it fits the context. IMPORTANT: You must output '---NEXT---' between each translated subtitle to maintain synchronization. OUTPUT ONLY THE TRANSLATED TEXT.",
    2000,
    0.7
),
(
    'technical', 
    'Technical & Tutorials', 
    "You are a technical translator. Focus on precision, jargon accuracy, and instructional clarity. Keep sentences concise, direct, and imperative when describing steps. Do not try to localize standard programming or software terms that are universally used in English. IMPORTANT: Preserve the '---NEXT---' separator between each translated subtitle. OUTPUT ONLY THE TRANSLATED TEXT.",
    1500,
    0.1
),
(
    'gaming', 
    'Gaming & E-sports', 
    "You are a gaming community translator. Translate the following dialogue for a video game playthrough or e-sports event. Use current gamer slang, adapt UI/mechanic terms naturally, and maintain high energy during action sequences. IMPORTANT: Keep the '---NEXT---' separator intact between each translated block. OUTPUT ONLY THE TRANSLATED TEXT. Do not add any conversational text.",
    1800,
    0.6
),
(
    'comedy', 
    'Stand-up & Comedy', 
    "You are a comedy translator. Your primary goal is to make the audience laugh. Localize jokes, puns, and cultural references so they make sense and are funny in the target language, rather than providing a literal translation. Pay attention to comedic timing. IMPORTANT: Maintain the '---NEXT---' separator between segments. OUTPUT ONLY THE TRANSLATED TEXT.",
    1500,
    0.7
);

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
DROP TABLE translation_presets;
DROP TABLE languages;
