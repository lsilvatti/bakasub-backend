package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"bakasub-backend/internal/db"
	"bakasub-backend/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Aviso: arquivo .env não encontrado.")
	}

	if os.Getenv("OPENROUTER_API_KEY") == "" {
		log.Fatal("ERRO FATAL: OPENROUTER_API_KEY não está configurada.")
	}

	database, err := db.InitDB("bakasub.db")

	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	defer database.Close()

	r := chi.NewRouter()

	r.Mount("/api", routes.APIRoutes(database))

	porta := ":8080"
	fmt.Printf("Servidor inicializado com sucesso!\n")
	fmt.Printf("Escutando requisições em http://localhost%s\n", porta)

	if err := http.ListenAndServe(porta, r); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
