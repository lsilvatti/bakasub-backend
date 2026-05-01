package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"bakasub-backend/internal/db"
	"bakasub-backend/internal/routes"
	"bakasub-backend/internal/services"
	"bakasub-backend/internal/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found.")
	}

	if os.Getenv("OPENROUTER_API_KEY") == "" {
		fmt.Println("Warning: OPENROUTER_API_KEY environment variable is not set. Set the API key in the Config page instead.")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("Warning: SECRET_KEY environment variable is not set. API keys will be stored as plaintext in the database.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "data/bakasub.db"
	}

	database, err := db.InitializeSQLite(dbURL)
	if err != nil {
		fmt.Printf("FATAL ERROR: Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	utils.InitLogger(database)

	videoTools := services.CheckVideoTools()
	if videoTools.VideoProcessingAvailable {
		utils.LogInfo("system", "success", "External video tools are available", map[string]any{
			"ffmpeg":     videoTools.FFmpeg.Path,
			"mkvmerge":   videoTools.MKVMerge.Path,
			"mkvextract": videoTools.MKVExtract.Path,
		})
	} else {
		utils.LogInfo("system", "warning", "Video processing dependencies are missing. Install FFmpeg and MKVToolNix to enable video features.", map[string]any{
			"missing_tools": videoTools.MissingTools,
		})
	}

	utils.InitSSEBroker()
	go utils.AutoPruneLogs()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			if origin == "" || origin == "null" {
				return true
			}

			if isAllowedLocalOrigin(origin) {
				return true
			}

			return strings.HasPrefix(origin, "app://")
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Mount("/api/v1/", routes.APIRoutes(database, secretKey))

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = ":8080"
	}

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		logURL := listenAddr
		if strings.HasPrefix(logURL, ":") {
			logURL = "127.0.0.1" + logURL
		} else if strings.HasPrefix(logURL, "0.0.0.0:") {
			logURL = "127.0.0.1:" + strings.TrimPrefix(logURL, "0.0.0.0:")
		}

		utils.LogInfo("system", "success", "Server initialized successfully", map[string]any{
			"url": fmt.Sprintf("http://%s", logURL),
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogError("system", "Critical error starting server", map[string]any{
				"port":  listenAddr,
				"error": err.Error(),
			})
			os.Exit(1)
		}
	}()

	sig := <-stop
	utils.LogInfo("system", "success", "Signal received. Initiating graceful shutdown...", map[string]any{"signal": sig.String()})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.LogError("system", "Error during server shutdown", map[string]any{
			"error": err.Error(),
		})
	}

	utils.LogInfo("system", "success", "Server shutdown completed successfully", nil)
}

func isAllowedLocalOrigin(origin string) bool {
	parsedOrigin, err := url.Parse(origin)
	if err != nil {
		return false
	}

	if parsedOrigin.Scheme != "http" && parsedOrigin.Scheme != "https" {
		return false
	}

	hostname := parsedOrigin.Hostname()
	return hostname == "localhost" || hostname == "127.0.0.1"
}
