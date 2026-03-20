package routes

import (
	"bakasub-backend/internal/ai"
	"bakasub-backend/internal/fileio"
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"

	"github.com/go-chi/chi/v5"
)

func TranslateRoutes() chi.Router {
	r := chi.NewRouter()

	openRouterService := ai.NewOpenRouterService()
	diskService := fileio.NewFileIOService()

	translationService := services.NewTranslatorService(openRouterService, diskService)

	translateHandler := &handlers.TranslateHandler{
		Translator: translationService,
	}

	r.Post("/", translateHandler.Translate)

	return r
}
