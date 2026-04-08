package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"bakasub-backend/internal/db"
	"bakasub-backend/internal/routes"
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
		dbURL = "postgres://bakasub:bakasub_password@localhost:5432/bakasub?sslmode=disable"
	}

	database, err := db.InitializePostgres(dbURL)
	if err != nil {
		fmt.Printf("FATAL ERROR: Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	utils.InitLogger(database)

	utils.InitSSEBroker()
	go utils.AutoPruneLogs()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Mount("/api/v1/", routes.APIRoutes(database, secretKey))

	port := ":8080"

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		utils.LogInfo("system", "success", "Server initialized successfully", map[string]any{
			"url": fmt.Sprintf("http://localhost%s", port),
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogError("system", "Critical error starting server", map[string]any{
				"port":  port,
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
