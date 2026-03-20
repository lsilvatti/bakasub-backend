# Bakasub - AI Subtitle Translator

Subtitle Translate - Powered by AI

> *"It's not like I made a tool for you to translate your subtitles, baka!"* - Bakasub

Bakasub is a self-hosted tool that automatically extracts, translates using Artificial Intelligence, and remuxes subtitles directly into your video files (MKV, MP4) using the LLM of your choice.

## 🚀 How to Install (Using Docker)

1. **Clone this repository:**
   ```bash
   git clone https://github.com/lsilvatti/bakasub.git
   cd bakasub
   ```

2. **Configure your environment variables:**
   ```bash
   cp .env.example .env
   ```
   Open the `.env` file and insert your `OPENROUTER_API_KEY` and the absolute path to your local videos directory (`VIDEO_DIR`).

3. **Start the application:**
   ```bash
   docker compose up -d --build
   ```

4. **Access the dashboard:**
   Open your browser and navigate to: **http://localhost:3000**

## 🛠️ Architecture
* **Backend:** Go (Golang) + SQLite
* **Frontend:** React + TypeScript + Vite
* **Video Processing:** FFmpeg + MKVToolNix
* **Orchestration:** Docker Compose