# Bakasub - AI Subtitle Translator 🌸✨

Subtitle Translate - Powered by AI.

> *"I-It's not like I made a tool for you to translate your subtitles, baka!"* Listen up! This is the **main repository**. You don't need to look at the messy backend or frontend code if you just want to use the app. Bakasub is a self-hosted tool that automatically extracts, translates using top-tier AI models (via OpenRouter), and remuxes subtitles directly into your video files (MKV, MP4) without breaking a sweat.

I orchestrate everything for you using Docker, building the backend and frontend straight from their main branches. So just follow the instructions and don't mess it up!

## 🛠️ What I'm Made Of (Architecture)
* **Backend:** Go (Golang) + PostgreSQL (I upgraded from SQLite, try to keep up!)
* **Frontend:** React + TypeScript + Vite (served via Nginx)
* **Video Processing:** FFmpeg + MKVToolNix
* **Orchestration:** Docker Compose

## 🚀 How to Date Me (Installation)

1. **Clone this repository:**
   ```bash
   git clone [https://github.com/lsilvatti/bakasub.git](https://github.com/lsilvatti/bakasub.git)
   cd bakasub
   ```

2. **Configure your environment variables:**
   ```bash
   cp .env.example .env
   ```
   Open the `.env` file. You *must* insert:
   * Your `OPENROUTER_API_KEY`
   * Your `TMDB_API_KEY` (So I can fetch movie metadata for you. Not that I care about what you watch!)
   * The absolute path to your local videos directory (`VIDEO_DIR`). If you don't give me a video directory, I won't have anything to translate!

3. **Start the application:**
   ```bash
   docker compose up -d --build
   ```
   *Note: I will automatically spin up a PostgreSQL 16 database, download the latest backend/frontend code, and wire everything together. Be patient while I build!*

4. **Access the dashboard:**
   Open your browser and navigate to: **http://localhost:3000**