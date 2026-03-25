package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func LanguageRoutes(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	languageService := services.NewLanguageService(db)
	languageHandler := &handlers.LanguageHandler{Service: languageService}

	r.Get("/", languageHandler.GetLanguages)
	r.Post("/", languageHandler.CreateLanguage)
	r.Put("/", languageHandler.UpdateLanguage)
	r.Delete("/", languageHandler.DeleteLanguage)

	return r
}
