package routes

import (
	"bakasub-backend/internal/handlers"
	"bakasub-backend/internal/services"
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func VideoRoutes(db *sql.DB, secretKey string) chi.Router {
	r := chi.NewRouter()

	videoService := services.NewVideoService()
	configService := services.NewConfigService(db, secretKey)

	videoHandler := handlers.VideoHandler{
		Processor: videoService,
		Config:    configService,
	}

	r.Get("/get-tracks", videoHandler.GetTrackHandler)
	r.Post("/extract-track", videoHandler.ExtractTrackHandler)
	r.Post("/merge-subtitle", videoHandler.MergeTrackHandler)

	return r
}
