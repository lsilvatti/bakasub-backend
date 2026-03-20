package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func ConfigRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	configService := services.NewConfigService(database)

	configHandler := &handlers.ConfigHandler{
		Service: configService,
	}

	r.Get("/", configHandler.GetUserConfig)
	r.Put("/", configHandler.UpdateUserConfig)

	return r
}
