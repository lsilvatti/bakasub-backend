package main

import (
	"context"
	"fmt"
	"log"
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
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found.")
	}

	if os.Getenv("OPENROUTER_API_KEY") == "" {
		log.Fatal("FATAL ERROR: OPENROUTER_API_KEY is not set.")
	}

	database, err := db.InitDB("bakasub.db")

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer database.Close()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Mount("/api", routes.APIRoutes(database))

	port := ":8080"

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		fmt.Printf("Server initialized successfully!\n")
		fmt.Printf("Listening for requests at http://localhost%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	sig := <-stop
	fmt.Printf("\nSignal received (%v). Initiating graceful shutdown...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	fmt.Println("Server successfully shut down.")
}
