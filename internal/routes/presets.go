package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func PresetRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	PresetService := services.NewPresetService(database)

	PresetHandler := &handlers.PresetHandler{
		Service: PresetService,
	}

	r.Get("/", PresetHandler.GetPresetsHandler)
	r.Post("/", PresetHandler.CreatePresetHandler)
	r.Put("/", PresetHandler.UpdatePresetHandler)
	r.Delete("/", PresetHandler.DeletePresetHandler)

	return r
}
