# Bakasub - API & Core Engine

This is the backend repository for **Bakasub**, responsible for all video processing logic, subtitle extraction, LLM (Artificial Intelligence) communication, and remuxing.

> **Note:** This repository is part of a larger ecosystem. If you just want to install and use Bakasub, check out the [main orchestration repository](https://github.com/lsilvatti/bakasub).

## 🛠️ Technologies Used
* **Language:** Go (Golang) 1.22+
* **Router:** `go-chi/chi`
* **Database:** SQLite (`modernc.org/sqlite` - pure Go driver without CGO)
* **Migrations:** Goose
* **Media Processing:** FFmpeg & MKVToolNix

## 🏗️ Architecture

The backend was built using **Clean Architecture** and **Dependency Inversion (SOLID)** principles. 
The business logic (Services) is completely decoupled from the transport layer (HTTP/Handlers) and external tools (Disk, Database, AI APIs).

* `internal/handlers`: Handles only HTTP requests and responses.
* `internal/services`: Pure business rules (grouped by domain).
* `internal/ai` & `internal/fileio`: Concrete implementations of the external world (OpenRouter, OS/Disk).
* Everything communicates via **Interfaces**, making the code 100% testable.

## 💻 Local Development

If you wish to contribute or modify the API locally:

### Prerequisites
Ensure you have Go, FFmpeg, and MKVToolNix installed on your system (e.g., `sudo pacman -S go ffmpeg mkvtoolnix-cli`).

### Running the project
For the best development experience, we recommend using [Air](https://github.com/cosmtrek/air) for live-reloading:

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Start the server with Air:
   ```bash
   air
   ```
   *The API will be running at `http://localhost:8080`.*