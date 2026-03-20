package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func LanguageRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	LanguageService := services.NewLanguageService(database)

	LanguageHandler := &handlers.LanguageHandler{
		Service: LanguageService,
	}

	r.Get("/", LanguageHandler.GetLanguagesHandler)
	r.Post("/", LanguageHandler.AddLanguageHandler)
	r.Put("/", LanguageHandler.UpdateLanguageHandler)
	r.Delete("/", LanguageHandler.DeleteLanguageHandler)

	return r
}
