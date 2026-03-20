package routes

import (
	"bakasub-backend/internal/ai"
	"bakasub-backend/internal/fileio"
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func OpenRouterTranslateRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	openRouterService := ai.NewOpenRouterService()
	diskService := fileio.NewFileIOService()

	translationService := services.NewTranslatorService(openRouterService, diskService, database)

	translateHandler := &handlers.TranslateHandler{
		Translator: translationService,
	}

	r.Post("/translate", translateHandler.Translate)

	return r
}
