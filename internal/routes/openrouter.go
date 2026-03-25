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
	jobService := services.NewJobService(database)

	translationService := services.NewTranslatorService(openRouterService, diskService, database, jobService)

	translateHandler := &handlers.TranslateHandler{
		Translator: translationService,
		JobService: jobService,
	}

	r.Post("/translate", translateHandler.Translate)
	r.Post("/preflight", translateHandler.PreFlight)

	return r
}
