package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func PresetRoutes(database *sql.DB) chi.Router {
	r := chi.NewRouter()

	presetService := services.NewPresetService(database)

	presetHandler := &handlers.PresetHandler{
		Service: presetService,
	}

	r.Get("/", presetHandler.GetPresets)
	r.Post("/", presetHandler.CreatePreset)
	r.Put("/", presetHandler.UpdatePreset)
	r.Delete("/", presetHandler.DeletePreset)

	return r
}
