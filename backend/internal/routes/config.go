package routes

import (
	"bakasub-backend/internal/handlers"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

// O retorno chi.Router é a melhor prática aqui
func ConfigRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()
	configHandler := handlers.ConfigHandler{DB: database}

	r.Get("/", configHandler.GetUserConfig)
	r.Put("/", configHandler.UpdateUserConfig)

	return r
}
