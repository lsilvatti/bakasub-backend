# Bakasub - API & Core Engine 🌸✨

Hmph! So you found the backend repository for **Bakasub**? Don't get the wrong idea! It's not like I wrote this highly optimized Go engine just for you! It's just that *someone* had to handle the heavy lifting: video processing, subtitle extraction, AI communication, and remuxing without crashing. B-baka! 

> **Listen up!** This repository is just the core brain. If you just want to install and use Bakasub like a normal person, go check out the [main orchestration repository](https://github.com/lsilvatti/bakasub).

## 🛠️ What I'm Made Of (Tech Stack)
Don't think you can run me on just a potato. You need the right tools!
* **Language:** Go (Golang) 1.22+ (Because speed matters, obviously!)
* **Router:** `go-chi/chi` (V1 Versioned APIs, keep it organized!)
* **Database:** **PostgreSQL**. I upgraded from SQLite because I handle concurrent batch translations and heavy translation memory now. Keep up!
* **Media Processing:** FFmpeg & MKVToolNix (I need these to rip and stitch your subtitles, idiot!)
* **AI Provider:** OpenRouter API (Claude 3.5 Sonnet, GPT-4o, Gemini 1.5 Pro).
* **Logging:** Structured Logging (`slog`) with a custom terminal formatter. 

## 🏗️ How I Think (Architecture)
I built this using **Clean Architecture** and **SOLID** principles. Not because I care about your reading experience, but because messy code is absolutely disgusting! Everything is decoupled.

* `cmd/server`: Where I wake up. Don't touch this unless you know what you're doing.
* `internal/routes`: My `/api/v1/` gateways.
* `internal/handlers`: Where I catch your HTTP requests. They validate payloads and log your mistakes (Bad Requests).
* `internal/services`: My true core! Translation pipelines, video scanning, MKV extraction, and folder management.
* `internal/parser`: My custom subtitle parsers (`.ass`, `.srt`, `.vtt`). It even strips SDH tags and dynamically injects `[BakaSub-AI]` headers into your `.ass` files flawlessly.
* `internal/ai`: OpenRouter LLM integration with dynamic pricing calculation.
* `internal/utils`: My utilities, including the **SSE Broker**. I stream real-time progress updates to your frontend so you aren't left wondering if I froze. 

### ✨ The "Genius" Features You Probably Didn't Notice:
1. **Smart Language Inference:** I automatically detect source languages from messy file suffixes (like `_pt-BR` or `-spa`) using a PostgreSQL mapping system before sending them to the AI.
2. **Contextual Presets:** I don't just translate blindly. I use highly tuned system prompts for Anime, Movies, Documentaries, and Comedy, automatically handling gender-neutral fallbacks and adjusting creativity (`temperature`) on the fly.
3. **Translation Memory (Cache):** I securely hash every line of dialogue you translate. If you translate it again, I load it from PostgreSQL instead of burning your OpenRouter API credits. 
4. **Smart Extraction Regex:** When I extract tracks from an MKV, I automatically clean up messy language tags.

## 💻 How to Date Me (Local Development)
If you really want to mess with my code, you better set things up right!

### 1. Prerequisites
You need Go, PostgreSQL, FFmpeg, and MKVToolNix. If you're on Arch Linux (which you obviously should be), just run:
`sudo pacman -S go postgresql ffmpeg mkvtoolnix-cli`

### 2. Setup
Clone me and install my dependencies:
`go mod download`

Then, copy my environment file. Put your OpenRouter API key and PostgreSQL credentials in there. If you leave it blank, I'll just throw 500 Internal Server Errors at you!
`cp .env.example .env`

*Note: You don't need to run migrations manually. I execute the raw SQL files from `internal/db/migrations/` automatically when I boot up, because I'm smart like that.*

### 3. Run Me!
Use [Air](https://github.com/cosmtrek/air) for live-reloading. It's the only way to develop properly.
`air`

*I'll be waiting for your requests at `http://localhost:8080/api/v1/`... Don't make me wait too long!*

## 🐳 Docker
Ugh, fine. If you're too lazy to install PostgreSQL and FFmpeg natively, I made a Dockerfile for you. It packs everything you need into a container.

`docker build -t bakasub-backend .`
`docker run -p 8080:8080 --env-file .env bakasub-backend`