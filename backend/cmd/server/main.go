package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"bakasub-backend/internal/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Aviso: arquivo .env não encontrado.")
	}

	if os.Getenv("OPENROUTER_API_KEY") == "" {
		log.Fatal("ERRO FATAL: OPENROUTER_API_KEY não está configurada.")
	}

	http.HandleFunc("/api/translate", handlers.TranslateHandler)

	porta := ":8080"
	fmt.Printf("Servidor inicializado com sucesso!\n")
	fmt.Printf("Escutando requisições POST em http://localhost%s/api/translate\n", porta)

	if err := http.ListenAndServe(porta, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
